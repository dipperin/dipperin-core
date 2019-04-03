# Minemaster

## module analysis

![image](https://github.com/caiqingfeng/res/blob/master/06-hld/02-test-net/codes/master.png?raw=true)

1. Once the instance of `MasterServer`, `server` receives work submitted by 
`localConnector` message via `ReceiveMsg`, it verifies the work `difficulty`.
If this is a valid work, `server` submits the work by `s.onSubmitBlock(w)`.
`onSubmitBlock` gets the current block and seals it, then passes the sealed block
to `workManager`.
  
2. A default `workManager` is provided. Each `masterServer` has only one `workManager`
to manage all the workers' works. Each `worker` corresponds to one `workerPerformance` 
instance. Different master could utilize different performance metrics/indicator by
modifying their `workerPerformance`. Also, among the same `master`/`workManager`, 
different `workerPerformance` could by applied to distinct workers.

3. The default instance of `workManager`, `defaultWorkManager` receives the 
newly sealed block, it updates the performance of the specific worker by calling 
`manager.performance[workerAddress].updatePerformance()`.

4. `defaultPerformance` provides a simple performance indicator based on how many
blocks a worker has worked out. It increments the block count `blocksMined`, every
time `.updatePerformance()` is called.  

5. Performance is a reference for reward distribution. Reward distribution may vary
for different performance metrics and reward distribution policies. A default distribution
policy is given, which simply linear to the performance. Note that, reward distribution
is independent of how performance is calculated, it simply divide the coibase reward
according to the given performance. Here, we provided a fault distribution method, which
equally distributes the coinbase based on workers' performance.

## Design Problem

 - It is easy for us to keep record of how much work the workers have done. But it is
 difficult to actually work out their reward. Since the reward is unknown unless the 
 previous mined block is verified by `verifiers`, there is no way to determine reward
 distribution before Minemaster actually receives coinbase.
 
 - If the master decides to pay **salary** to workers constantly regardless of work
 submitted or how long they have join the pool, it should be straightforward to 
 distribute reward by performance. 
 
 ## Conclusion
 
 - It si enough to provide performance indicator and reward distribution template for 
 Minemasters, since it is not stated in the requirement documents.
