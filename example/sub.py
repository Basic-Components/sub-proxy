import sys
import argparse
import zmq


def _sub(args):
    context = zmq.Context()
    socket = context.socket(zmq.SUB)
    socket.connect(args.url)
    print(f"process subscribe @ {args.url} !")

    socket.setsockopt_string(zmq.SUBSCRIBE, args.topic)
    # Process 5 updates
    total_temp = 0
    for update_nbr in range(5):
        string = socket.recv_string()
        print(f"recv string {string}")
        topic, temperature, relhumidity = string.split(" ")
        total_temp += int(temperature)
    avg_t = total_temp / (update_nbr + 1)
    print(f"Average temperature for zipcode {topic} was {avg_t}F")


def _parser_args(params):
    """解析命令行参数."""
    parser = argparse.ArgumentParser()
    parser.add_argument('--url', type=str, default="tcp://localhost:5571", help="指定连接到哪个组件")
    parser.add_argument('--topic', type=str, default="10001", help="指定订阅的主题")
    parser.set_defaults(func=_sub)
    args = parser.parse_args(params)
    args.func(args)


def main(argv=sys.argv[1:]):
    """服务启动入口.

    设置覆盖顺序`命令行参数`>`'-c'指定的配置文件`>`项目启动位置的配置文件`>默认配置.
    """
    _parser_args(argv)


if __name__ == '__main__':
    try:
        main()
    except KeyboardInterrupt:
        print('- Ctrl+C pressed in Terminal')
        print("Server shutdown!")
