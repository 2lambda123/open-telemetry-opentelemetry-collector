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
	"math/rand"

	otlpcommon "github.com/open-telemetry/opentelemetry-proto/gen/go/common/v1"
	otlptrace "github.com/open-telemetry/opentelemetry-proto/gen/go/trace/v1"

	"github.com/open-telemetry/opentelemetry-collector/translator/conventions"
)

const (
	SpanAttrNil               = "Nil"
	SpanAttrEmpty             = "Empty"
	SpanAttrDatabaseSQL       = "DatabaseSQL"
	SpanAttrDatabaseNoSQL     = "DatabaseNoSQL"
	SpanAttrFaaSDatasource    = "FaaSDatasource"
	SpanAttrFaaSHTTP          = "FaaSHTTP"
	SpanAttrFaaSPubSub        = "FaaSPubSub"
	SpanAttrFaaSTimer         = "FaaSTimer"
	SpanAttrFaaSOther         = "FaaSOther"
	SpanAttrHTTPClient        = "HTTPClient"
	SpanAttrHTTPServer        = "HTTPServer"
	SpanAttrMessagingProducer = "MessagingProducer"
	SpanAttrMessagingConsumer = "MessagingConsumer"
	SpanAttrGRPCClient        = "gRPCClient"
	SpanAttrGRPCServer        = "gRPCServer"
	SpanAttrInternal          = "Internal"
)

func GenerateSpan(traceID []byte, parentID []byte, spanName string, kind string, spanTypeID string) *otlptrace.Span {
	return &otlptrace.Span{
		TraceId:                traceID,
		SpanId:                 generateSpanID(),
		TraceState:             "",
		ParentSpanId:           parentID,
		Name:                   spanName,
		Kind:                   otlptrace.Span_CLIENT,
		StartTimeUnixNano:      0,
		EndTimeUnixNano:        0,
		Attributes:             generateSpanAttributes(spanTypeID),
		DroppedAttributesCount: 0,
		Events:                 nil,
		DroppedEventsCount:     0,
		Links:                  nil,
		DroppedLinksCount:      0,
		Status: &otlptrace.Status{
			Code:    otlptrace.Status_Ok,
			Message: "",
		},
	}
}

func generateSpanID() []byte {
	var r [8]byte
	_, err := rand.Read(r[:])
	if err != nil {
		panic(err)
	}
	return r[:]
}

func generateSpanAttributes(spanTypeID string) []*otlpcommon.AttributeKeyValue {
	var attrs map[string]interface{}
	if SpanAttrNil == spanTypeID {
		attrs = nil
	} else if SpanAttrEmpty == spanTypeID {
		attrs = make(map[string]interface{})
	} else if SpanAttrDatabaseSQL == spanTypeID {
		attrs = generateDatabaseSQLAttributes()
	} else if SpanAttrDatabaseNoSQL == spanTypeID {
		attrs = generateDatabaseNoSQLAttributes()
	} else if SpanAttrFaaSDatasource == spanTypeID {
		attrs = generateFaaSDatasourceAttributes()
	} else if SpanAttrFaaSHTTP == spanTypeID {
		attrs = generateFaaSHTTPAttributes()
	} else if SpanAttrFaaSPubSub == spanTypeID {
		attrs = generateFaaSPubSubAttributes()
	} else if SpanAttrFaaSTimer == spanTypeID {
		attrs = generateFaaSTimerAttributes()
	} else if SpanAttrFaaSOther == spanTypeID {
		attrs = generateFaaSOtherAttributes()
	} else if SpanAttrHTTPClient == spanTypeID {
		attrs = generateHTTPClientAttributes()
	} else if SpanAttrHTTPServer == spanTypeID {
		attrs = generateHTTPServerAttributes()
	} else if SpanAttrMessagingProducer == spanTypeID {
		attrs = generateMessagingProducerAttributes()
	} else if SpanAttrMessagingConsumer == spanTypeID {
		attrs = generateMessagingConsumerAttributes()
	} else if SpanAttrGRPCClient == spanTypeID {
		attrs = generateGRPCClientAttributes()
	} else if SpanAttrGRPCServer == spanTypeID {
		attrs = generateGRPCServerAttributes()
	} else if SpanAttrInternal == spanTypeID {
		attrs = generateInternalAttributes()
	} else {
		panic("invalid spanTypeID")
	}
	return convertMapToAttributeKeyValues(attrs)
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
	attrMap[conventions.AttributeNetPeerPort] = 5432
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
	attrMap[conventions.AttributeNetPeerPort] = 443
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
	attrMap[conventions.AttributeHTTPStatusCode] = 201
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
	attrMap["processed.count"] = 256
	attrMap["processed.data"] = 14.46
	attrMap["processed.errors"] = false
	attrMap[conventions.AttributeEnduserID] = "unittest"
	return attrMap
}

func generateHTTPClientAttributes() map[string]interface{} {
	attrMap := make(map[string]interface{})
	attrMap[conventions.AttributeHTTPMethod] = "GET"
	attrMap[conventions.AttributeHTTPURL] = "https://opentelemetry.io/registry/"
	attrMap[conventions.AttributeHTTPStatusCode] = 200
	attrMap[conventions.AttributeHTTPStatusText] = "More Than OK"
	attrMap[conventions.AttributeEnduserID] = "unittest"
	return attrMap
}

func generateHTTPServerAttributes() map[string]interface{} {
	attrMap := make(map[string]interface{})
	attrMap[conventions.AttributeHTTPMethod] = "POST"
	attrMap[conventions.AttributeHTTPScheme] = "https"
	attrMap[conventions.AttributeHTTPServerName] = "api22.opentelemetry.io"
	attrMap[conventions.AttributeNetHostPort] = 443
	attrMap[conventions.AttributeHTTPTarget] = "/blog/posts"
	attrMap[conventions.AttributeHTTPFlavor] = "2"
	attrMap[conventions.AttributeHTTPStatusCode] = 201
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
	attrMap["messaging.sequence"] = 1
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
	attrMap[conventions.AttributeNetHostPort] = 8443
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
