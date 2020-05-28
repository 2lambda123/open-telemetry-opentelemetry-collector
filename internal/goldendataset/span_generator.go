// Copyright 2020, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package goldendataset

import (
	"fmt"
	"io"
	"time"

	otlpcommon "go.opentelemetry.io/collector/internal/data/opentelemetry-proto-gen/common/v1"
	otlptrace "go.opentelemetry.io/collector/internal/data/opentelemetry-proto-gen/trace/v1"
	"go.opentelemetry.io/collector/translator/conventions"
)

var statusCodeMap = map[PICTInputStatus]otlptrace.Status_StatusCode{
	SpanStatusOk:                 otlptrace.Status_Ok,
	SpanStatusCancelled:          otlptrace.Status_Cancelled,
	SpanStatusUnknownError:       otlptrace.Status_UnknownError,
	SpanStatusInvalidArgument:    otlptrace.Status_InvalidArgument,
	SpanStatusDeadlineExceeded:   otlptrace.Status_DeadlineExceeded,
	SpanStatusNotFound:           otlptrace.Status_NotFound,
	SpanStatusAlreadyExists:      otlptrace.Status_AlreadyExists,
	SpanStatusPermissionDenied:   otlptrace.Status_PermissionDenied,
	SpanStatusResourceExhausted:  otlptrace.Status_ResourceExhausted,
	SpanStatusFailedPrecondition: otlptrace.Status_FailedPrecondition,
	SpanStatusAborted:            otlptrace.Status_Aborted,
	SpanStatusOutOfRange:         otlptrace.Status_OutOfRange,
	SpanStatusUnimplemented:      otlptrace.Status_Unimplemented,
	SpanStatusInternalError:      otlptrace.Status_InternalError,
	SpanStatusUnavailable:        otlptrace.Status_Unavailable,
	SpanStatusDataLoss:           otlptrace.Status_DataLoss,
	SpanStatusUnauthenticated:    otlptrace.Status_Unauthenticated,
}

var statusMsgMap = map[PICTInputStatus]string{
	SpanStatusOk:                 "",
	SpanStatusCancelled:          "Cancellation received",
	SpanStatusUnknownError:       "",
	SpanStatusInvalidArgument:    "parameter is required",
	SpanStatusDeadlineExceeded:   "timed out after 30002 ms",
	SpanStatusNotFound:           "/dragons/RomanianLonghorn not found",
	SpanStatusAlreadyExists:      "/dragons/Drogon already exists",
	SpanStatusPermissionDenied:   "tlannister does not have write permission",
	SpanStatusResourceExhausted:  "ResourceExhausted",
	SpanStatusFailedPrecondition: "33a64df551425fcc55e4d42a148795d9f25f89d4 has been edited",
	SpanStatusAborted:            "",
	SpanStatusOutOfRange:         "Range Not Satisfiable",
	SpanStatusUnimplemented:      "Unimplemented",
	SpanStatusInternalError:      "java.lang.NullPointerException",
	SpanStatusUnavailable:        "RecommendationService is currently unavailable",
	SpanStatusDataLoss:           "",
	SpanStatusUnauthenticated:    "nstark is unknown user",
}

//GenerateSpans generates a slice of OTLP Span objects with the number of spans specified by the count input
//parameter. The startPos parameter specifies the line in the PICT tool-generated, test parameter
//combination records file specified by the pictFile parameter to start reading from. When the end record
//is reached it loops back to the first record. The random parameter injects the random number generator
//to use in generating IDs and other random values. Using a random number generator with the same seed value
//enables reproducible tests.
//
//The return values are the slice with the generated spans, the starting position for the next generation
//run and the error which caused the spans generation to fail. If err is not nil, the spans slice will
//have nil values.
func GenerateSpans(count int, startPos int, pictFile string, random io.Reader) ([]*otlptrace.Span, int, error) {
	pairsData, err := loadPictOutputFile(pictFile)
	if err != nil {
		return nil, 0, err
	}
	pairsTotal := len(pairsData)
	spanList := make([]*otlptrace.Span, count)
	index := startPos + 1
	var inputs []string
	var spanInputs *PICTSpanInputs
	var traceID []byte
	var parentID []byte
	for i := 0; i < count; i++ {
		if index >= pairsTotal {
			index = 1
		}
		inputs = pairsData[index]
		spanInputs = &PICTSpanInputs{
			Parent:     PICTInputParent(inputs[SpansColumnParent]),
			Tracestate: PICTInputTracestate(inputs[SpansColumnTracestate]),
			Kind:       PICTInputKind(inputs[SpansColumnKind]),
			Attributes: PICTInputAttributes(inputs[SpansColumnAttributes]),
			Events:     PICTInputSpanChild(inputs[SpansColumnEvents]),
			Links:      PICTInputSpanChild(inputs[SpansColumnLinks]),
			Status:     PICTInputStatus(inputs[SpansColumnStatus]),
		}
		switch spanInputs.Parent {
		case SpanParentRoot:
			traceID = generateTraceID(random)
			parentID = nil
		case SpanParentChild:
			// use existing if available
			if traceID == nil {
				traceID = generateTraceID(random)
			}
			if parentID == nil {
				parentID = generateSpanID(random)
			}
		}
		spanName := generateSpanName(spanInputs)
		spanList[i] = GenerateSpan(traceID, parentID, spanName, spanInputs, random)
		parentID = spanList[i].SpanId
		index++
	}
	return spanList, index, nil
}

func generateSpanName(spanInputs *PICTSpanInputs) string {
	return fmt.Sprintf("/%s/%s/%s/%s/%s/%s/%s", spanInputs.Parent, spanInputs.Tracestate, spanInputs.Kind,
		spanInputs.Attributes, spanInputs.Events, spanInputs.Links, spanInputs.Status)
}

//GenerateSpan generates a single OTLP Span based on the input values provided. They are:
//  traceID - the trace ID to use, should not be nil
//  parentID - the parent span ID or nil if it is a root span
//  spanName - the span name, should not be blank
//  spanInputs - the pairwise combination of field value variations for this span
//  random - the random number generator to use in generating ID values
//
//The generated span is returned.
func GenerateSpan(traceID []byte, parentID []byte, spanName string, spanInputs *PICTSpanInputs,
	random io.Reader) *otlptrace.Span {
	endTime := time.Now().Add(-50 * time.Microsecond)
	return &otlptrace.Span{
		TraceId:                traceID,
		SpanId:                 generateSpanID(random),
		TraceState:             generateTraceState(spanInputs.Tracestate),
		ParentSpanId:           parentID,
		Name:                   spanName,
		Kind:                   lookupSpanKind(spanInputs.Kind),
		StartTimeUnixNano:      uint64(endTime.Add(-215 * time.Millisecond).UnixNano()),
		EndTimeUnixNano:        uint64(endTime.UnixNano()),
		Attributes:             generateSpanAttributes(spanInputs.Attributes),
		DroppedAttributesCount: 0,
		Events:                 generateSpanEvents(spanInputs.Events),
		DroppedEventsCount:     0,
		Links:                  generateSpanLinks(spanInputs.Links, random),
		DroppedLinksCount:      0,
		Status:                 generateStatus(spanInputs.Status),
	}
}

func generateTraceState(tracestate PICTInputTracestate) string {
	switch tracestate {
	case TraceStateOne:
		return "lasterror=f39cd56cc44274fd5abd07ef1164246d10ce2955"
	case TraceStateFour:
		return "err@ck=80ee5638,rate@ck=1.62,rojo=00f067aa0ba902b7,congo=t61rcWkgMzE"
	case TraceStateEmpty:
		fallthrough
	default:
		return ""
	}
}

func lookupSpanKind(kind PICTInputKind) otlptrace.Span_SpanKind {
	switch kind {
	case SpanKindClient:
		return otlptrace.Span_CLIENT
	case SpanKindServer:
		return otlptrace.Span_SERVER
	case SpanKindProducer:
		return otlptrace.Span_PRODUCER
	case SpanKindConsumer:
		return otlptrace.Span_CONSUMER
	case SpanKindInternal:
		return otlptrace.Span_INTERNAL
	case SpanKindUnspecified:
		fallthrough
	default:
		return otlptrace.Span_SPAN_KIND_UNSPECIFIED
	}
}

func generateSpanAttributes(spanTypeID PICTInputAttributes) []*otlpcommon.AttributeKeyValue {
	var attrs map[string]interface{}
	switch spanTypeID {
	case SpanAttrNil:
		attrs = nil
	case SpanAttrEmpty:
		attrs = make(map[string]interface{})
	case SpanAttrDatabaseSQL:
		attrs = generateDatabaseSQLAttributes()
	case SpanAttrDatabaseNoSQL:
		attrs = generateDatabaseNoSQLAttributes()
	case SpanAttrFaaSDatasource:
		attrs = generateFaaSDatasourceAttributes()
	case SpanAttrFaaSHTTP:
		attrs = generateFaaSHTTPAttributes()
	case SpanAttrFaaSPubSub:
		attrs = generateFaaSPubSubAttributes()
	case SpanAttrFaaSTimer:
		attrs = generateFaaSTimerAttributes()
	case SpanAttrFaaSOther:
		attrs = generateFaaSOtherAttributes()
	case SpanAttrHTTPClient:
		attrs = generateHTTPClientAttributes()
	case SpanAttrHTTPServer:
		attrs = generateHTTPServerAttributes()
	case SpanAttrMessagingProducer:
		attrs = generateMessagingProducerAttributes()
	case SpanAttrMessagingConsumer:
		attrs = generateMessagingConsumerAttributes()
	case SpanAttrGRPCClient:
		attrs = generateGRPCClientAttributes()
	case SpanAttrGRPCServer:
		attrs = generateGRPCServerAttributes()
	case SpanAttrInternal:
		attrs = generateInternalAttributes()
	case SpanAttrMaxCount:
		attrs = generateMaxCountAttributes()
	default:
		attrs = generateGRPCClientAttributes()
	}
	return convertMapToAttributeKeyValues(attrs)
}

func generateStatus(statusStr PICTInputStatus) *otlptrace.Status {
	if SpanStatusNil == statusStr {
		return nil
	}
	return &otlptrace.Status{
		Code:    statusCodeMap[statusStr],
		Message: statusMsgMap[statusStr],
	}
}

func generateDatabaseSQLAttributes() map[string]interface{} {
	attrMap := make(map[string]interface{})
	attrMap[conventions.AttributeDBType] = "sql"
	attrMap[conventions.AttributeDBInstance] = "inventory"
	attrMap[conventions.AttributeDBStatement] =
		"SELECT c.product_catg_id, c.catg_name, c.description, c.html_frag, c.image_url, p.name FROM product_catg c OUTER JOIN product p ON c.product_catg_id=p.product_catg_id WHERE c.product_catg_id = :catgId"
	attrMap[conventions.AttributeDBUser] = "invsvc"
	attrMap[conventions.AttributeDBURL] = "jdbc:postgresql://invdev.cdsr3wfqepqo.us-east-1.rds.amazonaws.com:5432/inventory"
	attrMap[conventions.AttributeNetPeerIP] = "172.30.2.7"
	attrMap[conventions.AttributeNetPeerPort] = int64(5432)
	attrMap[conventions.AttributeEnduserID] = "unittest"
	return attrMap
}

func generateDatabaseNoSQLAttributes() map[string]interface{} {
	attrMap := make(map[string]interface{})
	attrMap[conventions.AttributeDBType] = "cosmosdb"
	attrMap[conventions.AttributeDBInstance] = "graphdb"
	attrMap[conventions.AttributeDBStatement] = "g.V().hasLabel('postive').has('age', gt(65)).values('geocode')"
	attrMap[conventions.AttributeDBURL] = "wss://contacttrace.gremlin.cosmos.azure.com:443/"
	attrMap[conventions.AttributeNetPeerIP] = "10.118.17.63"
	attrMap[conventions.AttributeNetPeerPort] = int64(443)
	attrMap[conventions.AttributeEnduserID] = "unittest"
	return attrMap
}

func generateFaaSDatasourceAttributes() map[string]interface{} {
	attrMap := make(map[string]interface{})
	attrMap[conventions.AttributeFaaSTrigger] = conventions.FaaSTriggerDataSource
	attrMap[conventions.AttributeFaaSExecution] = "DB85AF51-5E13-473D-8454-1E2D59415EAB"
	attrMap[conventions.AttributeFaaSDocumentCollection] = "faa-flight-delay-information-incoming"
	attrMap[conventions.AttributeFaaSDocumentOperation] = "insert"
	attrMap[conventions.AttributeFaaSDocumentTime] = "2020-05-09T19:50:06Z"
	attrMap[conventions.AttributeFaaSDocumentName] = "delays-20200509-13.csv"
	attrMap[conventions.AttributeEnduserID] = "unittest"
	return attrMap
}

func generateFaaSHTTPAttributes() map[string]interface{} {
	attrMap := make(map[string]interface{})
	attrMap[conventions.AttributeFaaSTrigger] = conventions.FaaSTriggerHTTP
	attrMap[conventions.AttributeHTTPMethod] = "POST"
	attrMap[conventions.AttributeHTTPScheme] = "https"
	attrMap[conventions.AttributeHTTPHost] = "api.opentelemetry.io"
	attrMap[conventions.AttributeHTTPTarget] = "/blog/posts"
	attrMap[conventions.AttributeHTTPFlavor] = "2"
	attrMap[conventions.AttributeHTTPStatusCode] = int64(201)
	attrMap[conventions.AttributeHTTPUserAgent] =
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.1 Safari/605.1.15"
	attrMap[conventions.AttributeEnduserID] = "unittest"
	return attrMap
}

func generateFaaSPubSubAttributes() map[string]interface{} {
	attrMap := make(map[string]interface{})
	attrMap[conventions.AttributeFaaSTrigger] = conventions.FaaSTriggerPubSub
	attrMap[conventions.AttributeMessagingSystem] = "sqs"
	attrMap[conventions.AttributeMessagingDestination] = "video-views-au"
	attrMap[conventions.AttributeMessagingOperation] = "process"
	attrMap[conventions.AttributeEnduserID] = "unittest"
	return attrMap
}

func generateFaaSTimerAttributes() map[string]interface{} {
	attrMap := make(map[string]interface{})
	attrMap[conventions.AttributeFaaSTrigger] = conventions.FaaSTriggerTimer
	attrMap[conventions.AttributeFaaSExecution] = "73103A4C-E22F-4493-BDE8-EAE5CAB37B50"
	attrMap[conventions.AttributeFaaSTime] = "2020-05-09T20:00:08Z"
	attrMap[conventions.AttributeFaaSCron] = "0/15 * * * *"
	attrMap[conventions.AttributeEnduserID] = "unittest"
	return attrMap
}

func generateFaaSOtherAttributes() map[string]interface{} {
	attrMap := make(map[string]interface{})
	attrMap[conventions.AttributeFaaSTrigger] = conventions.FaaSTriggerOther
	attrMap["processed.count"] = int64(256)
	attrMap["processed.data"] = 14.46
	attrMap["processed.errors"] = false
	attrMap[conventions.AttributeEnduserID] = "unittest"
	return attrMap
}

func generateHTTPClientAttributes() map[string]interface{} {
	attrMap := make(map[string]interface{})
	attrMap[conventions.AttributeHTTPMethod] = "GET"
	attrMap[conventions.AttributeHTTPURL] = "https://opentelemetry.io/registry/"
	attrMap[conventions.AttributeHTTPStatusCode] = int64(200)
	attrMap[conventions.AttributeHTTPStatusText] = "More Than OK"
	attrMap[conventions.AttributeEnduserID] = "unittest"
	return attrMap
}

func generateHTTPServerAttributes() map[string]interface{} {
	attrMap := make(map[string]interface{})
	attrMap[conventions.AttributeHTTPMethod] = "POST"
	attrMap[conventions.AttributeHTTPScheme] = "https"
	attrMap[conventions.AttributeHTTPServerName] = "api22.opentelemetry.io"
	attrMap[conventions.AttributeNetHostPort] = int64(443)
	attrMap[conventions.AttributeHTTPTarget] = "/blog/posts"
	attrMap[conventions.AttributeHTTPFlavor] = "2"
	attrMap[conventions.AttributeHTTPStatusCode] = int64(201)
	attrMap[conventions.AttributeHTTPUserAgent] =
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36"
	attrMap[conventions.AttributeHTTPRoute] = "/blog/posts"
	attrMap[conventions.AttributeHTTPClientIP] = "2001:506:71f0:16e::1"
	attrMap[conventions.AttributeEnduserID] = "unittest"
	return attrMap
}

func generateMessagingProducerAttributes() map[string]interface{} {
	attrMap := make(map[string]interface{})
	attrMap[conventions.AttributeMessagingSystem] = "nats"
	attrMap[conventions.AttributeMessagingDestination] = "time.us.east.atlanta"
	attrMap[conventions.AttributeMessagingDestinationKind] = "topic"
	attrMap[conventions.AttributeMessagingMessageID] = "AA7C5438-D93A-43C8-9961-55613204648F"
	attrMap["messaging.sequence"] = int64(1)
	attrMap[conventions.AttributeNetPeerIP] = "10.10.212.33"
	attrMap[conventions.AttributeEnduserID] = "unittest"
	return attrMap
}

func generateMessagingConsumerAttributes() map[string]interface{} {
	attrMap := make(map[string]interface{})
	attrMap[conventions.AttributeMessagingSystem] = "kafka"
	attrMap[conventions.AttributeMessagingDestination] = "infrastructure-events-zone1"
	attrMap[conventions.AttributeMessagingOperation] = "receive"
	attrMap[conventions.AttributeNetPeerIP] = "2600:1700:1f00:11c0:4de0:c223:a800:4e87"
	attrMap[conventions.AttributeEnduserID] = "unittest"
	return attrMap
}

func generateGRPCClientAttributes() map[string]interface{} {
	attrMap := make(map[string]interface{})
	attrMap[conventions.AttributeRPCService] = "PullRequestsService"
	attrMap[conventions.AttributeNetPeerIP] = "2600:1700:1f00:11c0:4de0:c223:a800:4e87"
	attrMap[conventions.AttributeNetHostPort] = int64(8443)
	attrMap[conventions.AttributeEnduserID] = "unittest"
	return attrMap
}

func generateGRPCServerAttributes() map[string]interface{} {
	attrMap := make(map[string]interface{})
	attrMap[conventions.AttributeRPCService] = "PullRequestsService"
	attrMap[conventions.AttributeNetPeerIP] = "192.168.1.70"
	attrMap[conventions.AttributeEnduserID] = "unittest"
	return attrMap
}

func generateInternalAttributes() map[string]interface{} {
	attrMap := make(map[string]interface{})
	attrMap["parameters"] = "account=7310,amount=1817.10"
	attrMap[conventions.AttributeEnduserID] = "unittest"
	return attrMap
}

func generateMaxCountAttributes() map[string]interface{} {
	attrMap := make(map[string]interface{})
	attrMap[conventions.AttributeHTTPMethod] = "POST"
	attrMap[conventions.AttributeHTTPScheme] = "https"
	attrMap[conventions.AttributeHTTPHost] = "api.opentelemetry.io"
	attrMap[conventions.AttributeNetHostName] = "api22.opentelemetry.io"
	attrMap[conventions.AttributeNetHostIP] = "2600:1700:1f00:11c0:1ced:afa5:fd88:9d48"
	attrMap[conventions.AttributeNetHostPort] = int64(443)
	attrMap[conventions.AttributeHTTPTarget] = "/blog/posts"
	attrMap[conventions.AttributeHTTPFlavor] = "2"
	attrMap[conventions.AttributeHTTPStatusCode] = int64(201)
	attrMap[conventions.AttributeHTTPStatusText] = "Created"
	attrMap[conventions.AttributeHTTPUserAgent] =
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36"
	attrMap[conventions.AttributeHTTPRoute] = "/blog/posts"
	attrMap[conventions.AttributeHTTPClientIP] = "2600:1700:1f00:11c0:1ced:afa5:fd77:9d01"
	attrMap[conventions.AttributeNetPeerIP] = "2600:1700:1f00:11c0:1ced:afa5:fd77:9ddc"
	attrMap[conventions.AttributeNetPeerPort] = int64(39111)
	attrMap["ai-sampler.weight"] = float64(0.07)
	attrMap["ai-sampler.absolute"] = false
	attrMap["ai-sampler.maxhops"] = int64(6)
	attrMap["application.create.location"] = "https://api.opentelemetry.io/blog/posts/806673B9-4F4D-4284-9635-3A3E3E3805BE"
	attrMap["application.svcmap"] = "Blogosphere"
	attrMap["application.abflags"] = "UIx=false,UI4=true,flow-alt3=false"
	attrMap["application.thread"] = "proc-pool-14"
	attrMap["application.session"] = "233CC260-63A8-4ACA-A1C1-8F97AB4A2C01"
	attrMap["application.persist.size"] = int64(1172184)
	attrMap["application.queue.size"] = int64(0)
	attrMap["application.validation.results"] = "Success"
	attrMap["application.job.id"] = "0E38800B-9C4C-484E-8F2B-C7864D854321"
	attrMap["application.service.sla"] = float64(0.34)
	attrMap["application.service.slo"] = float64(0.55)
	attrMap[conventions.AttributeEnduserID] = "unittest"
	attrMap[conventions.AttributeEnduserRole] = "poweruser"
	attrMap[conventions.AttributeEnduserScope] = "email profile administrator"
	return attrMap
}

func generateSpanEvents(eventCnt PICTInputSpanChild) []*otlptrace.Span_Event {
	if SpanChildCountNil == eventCnt {
		return nil
	}
	listSize := calculateListSize(eventCnt)
	eventList := make([]*otlptrace.Span_Event, listSize)
	for i := 0; i < listSize; i++ {
		eventList[i] = generateSpanEvent(i)
	}
	return eventList
}

func generateSpanLinks(linkCnt PICTInputSpanChild, random io.Reader) []*otlptrace.Span_Link {
	if SpanChildCountNil == linkCnt {
		return nil
	}
	listSize := calculateListSize(linkCnt)
	linkList := make([]*otlptrace.Span_Link, listSize)
	for i := 0; i < listSize; i++ {
		linkList[i] = generateSpanLink(i, random)
	}
	return linkList
}

func calculateListSize(listCnt PICTInputSpanChild) int {
	switch listCnt {
	case SpanChildCountOne:
		return 1
	case SpanChildCountTwo:
		return 2
	case SpanChildCountEight:
		return 8
	case SpanChildCountEmpty:
		fallthrough
	default:
		return 0
	}
}

func generateSpanEvent(index int) *otlptrace.Span_Event {
	t := time.Now().Add(-75 * time.Microsecond)
	return &otlptrace.Span_Event{
		TimeUnixNano:           uint64(t.UnixNano()),
		Name:                   "message",
		Attributes:             generateEventAttributes(index),
		DroppedAttributesCount: 0,
	}
}

func generateEventAttributes(index int) []*otlpcommon.AttributeKeyValue {
	attrMap := make(map[string]interface{})
	if index%2 == 0 {
		attrMap[conventions.AttributeMessageType] = "SENT"
	} else {
		attrMap[conventions.AttributeMessageType] = "RECEIVED"
	}
	attrMap[conventions.AttributeMessageID] = int64(index)
	attrMap[conventions.AttributeMessageCompressedSize] = int64(17 * index)
	attrMap[conventions.AttributeMessageUncompressedSize] = int64(24 * index)
	return convertMapToAttributeKeyValues(attrMap)
}

func generateSpanLink(index int, random io.Reader) *otlptrace.Span_Link {
	return &otlptrace.Span_Link{
		TraceId:                generateTraceID(random),
		SpanId:                 generateSpanID(random),
		TraceState:             "",
		Attributes:             generateLinkAttributes(index),
		DroppedAttributesCount: 0,
	}
}

func generateLinkAttributes(index int) []*otlpcommon.AttributeKeyValue {
	attrMap := generateMessagingConsumerAttributes()
	return convertMapToAttributeKeyValues(attrMap)
}
