# jssp
JavaScript Server Pages:four_leaf_clover::tulip::rose::hibiscus:

- v1.0 分支js引擎使用的是[robertkrimen/otto](https://github.com/robertkrimen/otto)
- v2.0 分支js引擎使用的是[dop251/goja](https://github.com/dop251/goja) ,并在v1.0的基础上完善了诸多功能,包括ts,引擎缓存,ast缓存,项目重构等等
- 但是由于golang实现的js引擎性能孱弱,并且在golang与v8之间没有更好的绑定库之前,本项目是不会更新或者接受pr的.
- 同时发现了通过nodejs的pkg打包生成的二进制文件同样可以做到当初本项目的设想,
- 所以最后决定使用nodejs来实现jssp [jssp-nodejs](https://github.com/javascript-server-page/jssp-nodejs)
