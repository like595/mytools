package vrabbitmqhelper

import (
	"fmt"
	"github.com/like595/mytools/vtools"
	"github.com/streadway/amqp"
	"time"
)

type MyRabbitMQ struct {
	url              string
	channel          *amqp.Channel
	listenerFunction ListenerFunction
}

/*
监听函数
*/
type ListenerFunction func()

/*
初始化Rabbitmq
*/
func (this *MyRabbitMQ) Init(rabbitURL string) bool {
	this.url = rabbitURL
	return this.connection()
}

/*
连接 Rabbitmq
*/
func (this *MyRabbitMQ) connection() bool {
	conn, err := amqp.Dial(this.url)

	if err != nil {
		vtools.SugarLogger.Error("初始化Rabbitmq失败，错误：", err)
		return false
	}

	ch, err := conn.Channel()
	if err != nil {
		vtools.SugarLogger.Error("初始化Rabbitmq失败，错误：", err)
		return false
	}
	this.channel = ch
	this.createExchange("pmLog.topic")
	return true
}

/*
设置监听
*/
func (this *MyRabbitMQ) SetListener(listenerFunction ListenerFunction) {
	this.listenerFunction = listenerFunction
}

/*
发送消息
*/
func (this *MyRabbitMQ) Publish(exchange, routingKey string, data *[]byte) error {

	// 检查交换机是否存在，如果不存在则创建交换机

	err := this.channel.Publish(
		// exchange - yours may be different
		exchange,
		routingKey,
		// mandatory - we don't care if there I no queue
		false,
		// immediate - we don't care if there is no consumer on the queue
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         *data,
			DeliveryMode: amqp.Persistent,
		})
	if err != nil {
		time.Sleep(time.Minute)
		this.connection()
		this.Publish(exchange, routingKey, data)
		this.listenerFunction()
	}
	return err
}

func (this *MyRabbitMQ) createExchange(exchange string) {
	// 获取所有交换机信息
	(*this.channel).ExchangeDeclare(exchange, "topic", true, false, false, true, nil)
}

/*
监听消息
exchangeName：交换机名称
queueName：队列名称
routingKey：key名称
handler：接收到消息回调函数
concurrency：并发
*/
func (this *MyRabbitMQ) StartConsumer(
	exchangeName,
	queueName,
	routingKey string,
	handler func(d amqp.Delivery) bool,
	concurrency int) error {

	paramArgs := make(map[string]interface{})
	paramArgs["x-message-ttl"] = 5000

	//创建队列
	_, err := this.channel.QueueDeclare(queueName,
		false,     // 是否持久化
		true,      // 是否自动删除(前提是至少有一个消费者连接到这个队列，之后所有与这个队列连接的消费者都断开时，才会自动删除。注意：生产者客户端创建这个队列，或者没有消费者客户端与这个队列连接时，都不会自动删除这个队列)
		false,     // 是否为排他队列（排他的队列仅对“首次”声明的conn可见[一个conn中的其他channel也能访问该队列]，conn结束后队列删除）
		false,     // 是否阻塞
		paramArgs) //自定义参数
	if err != nil {
		return err
	}

	//绑定队列、交换机、key
	err = this.channel.QueueBind(queueName, routingKey, exchangeName, false, nil)
	if err != nil {
		return err
	}

	// prefetch 4x as many messages as we can handle at once
	prefetchCount := concurrency * 4
	err = this.channel.Qos(prefetchCount, 0, false)
	if err != nil {
		return err
	}

	msgs, err := this.channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return err
	}

	// create a goroutine for the number of concurrent threads requested
	for i := 0; i < concurrency; i++ {
		//fmt.Printf("Processing messages on thread %v...\n", i)
		go func() {
			for msg := range msgs {
				// if tha handler returns true then ACK, else NACK
				// the message back into the rabbit queue for
				// another round of processing
				if handler(msg) {
					msg.Ack(false)
				} else {
					msg.Nack(false, true)
				}
			}
			fmt.Println("Rabbit consumer closed - critical Error")
			this.connection()
			this.listenerFunction()
			//os.Exit(1)
		}()
	}
	return nil
}

//
///*
//初始化连接
//*/
//func (conn MyRabbitMQ) getURL() string {
//
//	sql := `SELECT item_text,item_value
//		from sys_dict a,sys_dict_item b
//		where dict_code='websocketSetting' and a.id=b.dict_id
//		ORDER BY item_value`
//	mySqlDBHelper := vdbhelper.MySqlDBHelper{}
//	mySqlDBHelper.Open()
//	rows, err := mySqlDBHelper.Select(sql)
//	var leixing, val sql2.NullString
//	var ip, port, username, password string
//	if err == nil {
//		for rows.Next() {
//			rows.Scan(&val, &leixing)
//
//			switch leixing.String {
//			case "1":
//				ip = val.String
//				break
//			case "3":
//				port = val.String
//				break
//			case "4":
//				username = val.String
//				break
//			case "5":
//				password = val.String
//				break
//			}
//		}
//		rows.Close()
//	}
//	mySqlDBHelper.Close()
//	url := fmt.Sprintf("amqp://%s:%s@%s:%s", username, password, ip, port)
//	return url
//}
