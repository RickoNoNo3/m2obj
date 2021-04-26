此包中只有`FileSyncer`的测试，目的是为了打破循环引用。

`FileSyncer`的测试需要用到`Formatter`的实现，也就是`m2obj/m2json`；但是`m2obj/m2json`又引用了`m2obj`，因此直接把测试写在`m2obj`包里会造成`m2obj`和`m2obj/m2json`的循环引用，所以必须写出来。
