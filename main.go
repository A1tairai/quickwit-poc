package main

import (
    "encoding/json"
    "fmt"
    "math/rand"
    "sync"
    "sync/atomic"
    "os"
    "strconv"
    "time"
    "github.com/confluentinc/confluent-kafka-go/kafka"
)

type LogData struct {
    MsgId               string  `json:"msg_id"`
    AppDst              string  `json:"app_dst"`
    AppSrc              string  `json:"app_src"`
    AutoInstanceDst     string  `json:"auto_instance_dst"`
    AutoInstanceSrc     string  `json:"auto_instance_src"`
    AutoInstanceTypeDst string  `json:"auto_instance_type_dst"`
    AutoInstanceTypeSrc string  `json:"auto_instance_type_src"`
    AutoServiceDst      string  `json:"auto_service_dst"`
    AutoServiceSrc      string  `json:"auto_service_src"`
    AutoServiceTypeDst  string  `json:"auto_service_type_dst"`
    AzSrc               string  `json:"az_src"`
    BizSrc              string  `json:"biz_src"`
    ClientPort          string  `json:"client_port"`
    DbClusterName       string  `json:"db_cluster_name"`
    DbInstanceName      string  `json:"db_instance_name"`
    DbStatement         string  `json:"db_statement"`
    DfInfo              string  `json:"df_info"`
    Duration            float64 `json:"duration"`
    EpcDst              string  `json:"epc_dst"`
    EpcSrc              string  `json:"epc_src"`
    FeatureSrc          string  `json:"feature_src"`
    FlowId              string  `json:"flow_id"`
    FlowInfoId          string  `json:"flow_info_id"`
    IpDst               string  `json:"ip_dst"`
    IpSrc               string  `json:"ip_src"`
    IsTls               string  `json:"is_tls"`
    L3EpcDst            string  `json:"l3_epc_dst"`
    L3EpcSrc            string  `json:"l3_epc_src"`
    L7Protocol          string  `json:"l7_protocol"`
    MetricName          string  `json:"metric_name"`
    NetworkProtocol     string  `json:"network_protocol"`
    PodClusterSrc       string  `json:"pod_cluster_src"`
    PodIdDst            string  `json:"pod_id_dst"`
    PodIdSrc            string  `json:"pod_id_src"`
    PodNodeSrc          string  `json:"pod_node_src"`
    PodNsSrc            string  `json:"pod_ns_src"`
    PrivateIpDst        string  `json:"private_ip_dst"`
    PrivateIpSrc        string  `json:"private_ip_src"`
    ReqTcpSeq           string  `json:"req_tcp_seq"`
    RequestType         string  `json:"request_type"`
    RespTcpSeq          string  `json:"resp_tcp_seq"`
    ResponseCode        string  `json:"response_code"`
    ServerPort          string  `json:"server_port"`
    SignalSource        string  `json:"signal_source"`
    Status              string  `json:"status"`
    TapPortType         string  `json:"tap_port_type"`
    TapSide             string  `json:"tap_side"`
    TeamSrc             string  `json:"team_src"`
    Tenant              string  `json:"tenant"`
    Timestamp           string  `json:"timestamp"`
    Topic               string  `json:"topic"`
    Type                string  `json:"type"`
    VtapDst             string  `json:"vtap_dst"`
    VtapId              string  `json:"vtap_id"`
    VtapSrc             string  `json:"vtap_src"`
}

var logDataPool = sync.Pool{
    New: func() interface{} {
        return &LogData{}
    },
}

func randomIP() string {
    return fmt.Sprintf("%d.%d.%d.%d", rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256))
}

func getKubernetesClientId(prefix string) string {
    hostname := os.Getenv("HOSTNAME") // HOSTNAME 环境变量通常包含 Pod 名称
    if hostname == "" {
        hostname = "unknown"
    }
    return fmt.Sprintf("%s-%s", prefix, hostname)
}

// ProducerConfig 创建 Kafka 生产者的配置
func ProducerConfig(broker string) *kafka.ConfigMap {
    clientId := getKubernetesClientId("go-producer")
    return &kafka.ConfigMap{
        "bootstrap.servers": broker,
        "client.id":         clientId,
        "compression.type":  "snappy",
        "queue.buffering.max.messages": 2000,
        "queue.buffering.max.kbytes": 3000,
    }
}

// CreateProducer 创建并返回 Kafka 生产者
func CreateProducer(broker string) (*kafka.Producer, error) {
    producer, err := kafka.NewProducer(ProducerConfig(broker))
    if err != nil {
        return nil, err
    }
    return producer, nil
}

// SendJSONMessage 发送 JSON 消息到 Kafka
func SendJSONMessage(producer *kafka.Producer, topic string, message string) error {
    msg := &kafka.Message{
        TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
        Value:          []byte(message),
    }
    return producer.Produce(msg, nil)
}

var msgIdCounter int64

func generateLogData() *LogData {
    logData := logDataPool.Get().(*LogData)

    *logData = LogData{
        MsgId:               strconv.FormatInt(atomic.AddInt64(&msgIdCounter, 1), 10),
        AppDst:              "tf-bin-prod-quickwit",
        AppSrc:              "quickwit",
        AutoInstanceDst:     randomIP(),
        AutoInstanceSrc:     "quickwit-searcher-86987c5bfc-lb9q5",
        AutoInstanceTypeDst: "internet_ip",
        AutoInstanceTypeSrc: "pod",
        AutoServiceDst:      randomIP(),
        AutoServiceSrc:      "quickwit-searcher",
        AutoServiceTypeDst:  "internet_ip",
        AzSrc:               "ap-northeast-1d",
        BizSrc:              "searcher",
        ClientPort:          strconv.Itoa(rand.Intn(65535)),
        DbClusterName:       "tf-bin-prod-quickwit-searcher",
        DbInstanceName:      "tf-bin-prod-quickwit-searcher-0016-001",
        DbStatement:         "GET An architecture built for performance and scalability",
        DfInfo:              "l7_flow_log",
        Duration:            rand.Float64(),
        EpcDst:              "0",
        EpcSrc:              "0",
        FeatureSrc:          "official",
        FlowId:              strconv.FormatInt(rand.Int63(), 10),
        FlowInfoId:          strconv.FormatInt(rand.Int63(), 10),
        IpDst:               randomIP(),
        IpSrc:               randomIP(),
        IsTls:               strconv.Itoa(rand.Intn(2)),
        L3EpcDst:            "-2",
        L3EpcSrc:            "-2",
        L7Protocol:          "Redis",
        MetricName:          "ebpf_l7_redis",
        NetworkProtocol:     "TCP",
        PodClusterSrc:       "tf-bin-prod-be-eks",
        PodIdDst:            "0",
        PodIdSrc:            "293382",
        PodNodeSrc:          "ip-192-119-50-221.ap-east-10.compute.internal",
        PodNsSrc:            "backend",
        PrivateIpDst:        "true",
        PrivateIpSrc:        "true",
        ReqTcpSeq:           strconv.FormatInt(rand.Int63(), 10),
        RequestType:         "GET",
        RespTcpSeq:          strconv.FormatInt(rand.Int63(), 10),
        ResponseCode:        "-256",
        ServerPort:          "6379",
        SignalSource:        "Packet",
        Status:              "Success",
        TapPortType:         "Local NIC",
        TapSide:             "Other NIC",
        TeamSrc:             "backend",
        Tenant:              "",
        Timestamp:           time.Now().Format(time.RFC3339),
        Topic:               "ebpf_l7_redis",
        Type:                "session",
        VtapDst:             "ip-192-119-50-221.ap-east-10.compute.internal-V22128",
        VtapId:              strconv.FormatInt(rand.Int63(), 10),
        VtapSrc:             "ip-192-119-50-221.ap-east-10.compute.internal-V22128",
    }

    return logData
}

func main() {
    sendToKafka := os.Getenv("SEND_TO_KAFKA") == "true"
    broker := os.Getenv("KAFKA_BROKER")
    topic := os.Getenv("KAFKA_TOPIC")

    var producer *kafka.Producer
    var err error

    if sendToKafka {
        producer, err = CreateProducer(broker)
        if err != nil {
            fmt.Printf("Failed to create producer: %s\n", err)
            return
        }
        defer producer.Close()
    }

    // Start a Goroutine for Kafka events handling
    go func() {
        for e := range producer.Events() {
            switch ev := e.(type) {
            case *kafka.Message:
                if ev.TopicPartition.Error != nil {
                    fmt.Printf("Delivery failed: %v\n", ev.TopicPartition.Error)
                }
            }
        }
    }()

    numPerSecond := 600
    counter := 0

    generateTicker := time.NewTicker(time.Second / time.Duration(numPerSecond))
    defer generateTicker.Stop()

    countTicker := time.NewTicker(time.Second)
    defer countTicker.Stop()

    for {
        select {
        case <-generateTicker.C:
            logData := generateLogData()
            jsonData, err := json.Marshal(logData)
            if err != nil {
                fmt.Println("Error marshalling JSON:", err)
                continue
            }

            if sendToKafka {
                err := SendJSONMessage(producer, topic, string(jsonData))
                if err != nil {
                    fmt.Printf("Failed to send message: %s\n", err)
                }
            } else {
                fmt.Println(string(jsonData))
            }

            // Release the logData back to the pool
            logDataPool.Put(logData)

            counter++
        case <-countTicker.C:
            fmt.Printf("Tasks generated this second: %d\n", counter)
            counter = 0

            // Flush the producer to ensure all messages are sent
            if sendToKafka {
                producer.Flush(10)
            }
        }
    }
}

