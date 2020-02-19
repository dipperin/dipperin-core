# mine worker

## module analysis

![image](https://github.com/caiqingfeng/res/blob/master/06-hld/02-test-net/codes/mine_worker.png?raw=true)

Above is the realization of mine worker. The main ideas are as follows:

1. Only a few functions such as start, stop, setting coinbase, consult coinbase are exposed to the external. The red box in the figure is the external function entry, so the other modules are not visible to the outside except the Worker is visible to the outside.

2. The internal implementation completely isolates the most likely changes in the future, and the orange part of the picture is the most likely place to change. According to the open-closed principle, the implemented code should not be modified for function, but should be modified or added by an extended method. Separating changes is a very important part.

3. The most likely changes to the mine worker are the mining algorithm and the received task data. The second thing that may change is the way of communication. Therefore, these two parts need to be isolated to deal with future changes. Other modules are called by the interface when calling these modules that may be changed. Therefore, only replace the used instance when replacing, do not modify the previous code.

4. The important modules in the mine worker are those in the above figure, and the various interfaces and various function types declared in the source file are made to isolate changes and module decoupling.

5. After decoupling by such a module, in addition to complying with the open-closed principle, and because these modules are small and more in line with the single responsibility principle, it is also convenient to write unit tests.

## Change method

1. If you want to modify the communication method, you can use the new communication method in the NewWorker method assembly method of mineworker/worker.go. Do not modify the original communication method.

2. If you want to modify the mining algorithm, add the corresponding task data in minemsg/messages.go, then add a new mining algorithm in mineworker/work_executor.go, and finally in mineworker/executor_builder.go Add a new method to build the executor, note that you should not directly modify the original mining algorithm.
