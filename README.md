# Demo Project

## Introduction
For Proof-of-Work (PoW) protection against DDoS attacks, we need to force an attacker to expend 
significantly more computational effort than what the server requires to verify a solution. 
At the same time, this mechanism must not degrade the user experience for legitimate users. 
One effective approach is to use a PoW algorithm similar to Bitcoin’s, where validity is 
determined solely by counting the number of leading zero bits in the hash.

### Advantages:
* Efficient Verification: Minimal computational overhead for the server.
* Asymmetric Workload: Clients must perform heavy computations while the server only does a couple of hash operations.
* Simple & Adjustable: Easily tune the difficulty by changing the required number of leading zero bits.
* Proven Security: Inspired by Bitcoin’s robust PoW mechanism, ensuring reliable protection.
* Broad Platform Support: SHA256 is available on most popular platforms and languages, making it easy to implement clients in various environments.

## Requirements
* Docker version 27.5.1

## Running the project

```shell
make run
```