### 卫星模块
***
server模块，route模块，client模块，dialer模块是代理运行的核心，我们可以
把他看作一个整体。而satellite（卫星）模块则是运行在这个主体之外的一个模块，
它的设定是围绕在核心旁辅助核心模块来完成更多能力，如：异步日志，上报埋点等。
犹如卫星一般