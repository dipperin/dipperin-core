## Gas Calculation

For normal transaction and contract transaction, the calculation of gas is slightly different:
### Normal Transaction

The total `gasUsed` value of a normal transaction is divided into two parts, `fixedGasUsed` is built into the system, and `f(txExtraData)` is calculated based on the `extraData` inside the transaction:

​                  `totalGasUsed = fixedGasUsed + f(txExtraData)`

Suppose that in a normal transaction, the number of bytes whose value is 0 in `exteraData` is `ZeroBytes`, and the number of bytes whose value is non-zero is `NoZeroBytes`. Then `f(txExtraData)` is expressed as follows:

​                `f(txExtraData) = TxDataZeroGas*ZeroBytes + TxDataNonZeroGas*NoZeroBytes`

| Parameters       | System Default Value | Remarks                                                      |
| ---------------- | -------------------- | ------------------------------------------------------------ |
| fixedGasUsed     | 21000                | This is the default `fixedGasUsed` value for the current system normal transaction. |
| TxDataZeroGas    | 4                    | When the data is 0, the `gasUsed` of unit bytes              |
| TxDataNonZeroGas | 68                   | When the data is not 0, the `gasUsed` of unit bytes          |

When the data is 0 and non-zero, its gasUsed is different because non-zero data consumes less system resources during storage and calculation.