# Demo Project

## Introduction
For Proof-of-Work (PoW) protection against DDoS attacks, we need to force an attacker to expend 
significantly more computational effort than what the server requires to verify a solution. 
At the same time, this mechanism must not degrade the user experience for legitimate users. 
One effective approach is to use a PoW algorithm similar to Bitcoin’s, where validity is 
determined solely by counting the number of leading zero bits in the hash.

### Advantages:
* Efficient Verification: Minimal computational overhead for the server.
* Asymmetric Workload: Clients (or attackers) must do heavy work, while verification is fast.
* Simple & Adjustable: Easily tune difficulty by changing the required leading zero bits.
* Proven Security: Inspired by Bitcoin’s robust PoW mechanism, ensuring reliable protection.

## Requirements
* Docker version 27.5.1

## Running the project

```shell
make run
```