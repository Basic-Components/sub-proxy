import sys
import argparse
import zmq
from random import randrange
import time


def _pub(args):
    #  Prepare our context and sockets
    context = zmq.Context()
    socket = context.socket(zmq.PUB)
    socket.bind(args.url)
    print(f"process pub @ {args.url} !")

    while True:
        topic = 10001
        temperature = randrange(-80, 135)
        relhumidity = randrange(10, 60)
        socket.send_string("%i %i %i" % (topic, temperature, relhumidity))
        print(f"temperature: {temperature},relhumidity {relhumidity} published @ topic {topic}")
        time.sleep(1)


def _parser_args(params):
    """解析命令行参数."""
    parser = argparse.ArgumentParser()
    parser.add_argument('--url', type=str, default="tcp://*:5570", help="指定连接到哪个组件")
    parser.set_defaults(func=_pub)
    args = parser.parse_args(params)
    args.func(args)


def main(argv=sys.argv[1:]):
    u"""服务启动入口.

    设置覆盖顺序`命令行参数`>`'-c'指定的配置文件`>`项目启动位置的配置文件`>默认配置.
    """
    _parser_args(argv)


if __name__ == '__main__':
    main()
