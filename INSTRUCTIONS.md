# SeatGeek's Coding Challenge

At SeatGeek, we do not believe in bringing engineers to our offices for long and intimidating sessions of brain-teasers and "whiteboard coding" interviews. Instead, we have found that the best way to assess an engineer's skill is to get them to work on a problem and situation similar to what they would find on their day-by-day working here.

The goal of this coding challenge is to have you produce code that shows us in concrete terms how you think about software engineering in your professional life. We want you to use the languages, tools, and setup with which you feel most comfortable. 

## Functional requirements
For this challenge, we ask you to write a service that manages the inventory component of a fictional distributed ticketing system. This scenario is an extremely simplified version of some of the challenges that our engineering teams face at SeatGeek.

For this exercise, a *Seat* is an inventory item that we sell to customers. For this exercise, we won't need to worry about events, venues, seat maps, or any other real-world complications. All we care about is seats that have been reserved or sold through the inventory system you will build.

The server will open a TCP socket on port 8099. It must accept connections from multiple clients at the same time, all trying issuing commands concurrently. Clients are independent of each other, and you should expect them to send repeated or contradicting messages. New clients can connect and disconnect at any moment, and sometimes clients can misbehave and send broken messages. Messages from clients follow this pattern:

```
<VERB>: <SEAT>\n
```

Where:
* `<VERB>` is either `RESERVE`, `BUY`, or `QUERY`
* `<SEAT>`  is an alphanumeric string identifying the seat referred to by the command, e.g. `A1`, `B12a`, `Back32`, etc.
* Both components are mandatory
* The message always ends with the character `\n`

For each message sent, the client will wait for a response code from the server. Expected responses are:

* For `RESERVE` commands, the service must return `OK` if the seat was previously `FREE`. If the seat was in any other state, the response should be `FAIL`
* For `BUY` commands, the service must return `OK` if the service was previously marked as `RESERVED`. If the seat was in any other state, the response should be `FAIL`
: The command was successfully  executed
* For `QUERY` commands, the service should return `FREE` if the queried seat hasn't been previously reserved or bought,  `RESERVED` if the seat has been reserved but not yet bought, and `SOLD` when the seat has already been bought
* The service should return `FAIL` for any unknown or invalid message it receives

### Non functional requirements
You should feel free to write your solution in any programming language in which you have professional experience in writing production-ready code. 

We would like to see as much code written by you as possible. That's why **we ask that you do not use any third-party software (e.g., libraries) apart from your chosen programming language's standard library**. If you absolutely must use some third party code, please write a sentence or two describing your motivation.

One typical exception to the rule above is for build tools and testing libraries, as these aren't part of most language's standard libraries and tools.

We also ask you to write code that you would consider production-ready, code that you would be comfortable deploying to production and maintaining.

Please don't forget to send us your automated tests and any other artifact needed to develop, build, or run your solution.

## The package we sent you

Together with this `INSTRUCTIONS.md` file, you should have received a tarball. In this tarball you will find:

* Executable files containing the test harness described below
* Another tarball, containing the Go source code for the executable mentioned above

### The test harness

So that you can verify the correctness of your solution before submitting it to us, we provide you with an automated test harness. 

To run the test suite, make sure your server is up and listening on port `8099`. Then execute the following command:

```
$ ./seatgeek-be-challenge
```

The tool will perform various checks to test your solution's robustness and correctness.  You should make sure that all tests pass before submitting your solution to the challenge. When all tests pass, you will see a message like this:

```
--------------------------
 ✅ TEST SUCCESSFUL ✅
--------------------------
```

We have built several other features into the test suite you might find helpful. To see them all, execute the following command:

```
$ ./seatgeek-be-challenge -help
```

## What do we expect from your submission

### No personally identifiable information

We are an equal opportunity employer and value diversity at our company. We do not discriminate on the basis of race, religion, color, national origin, gender, sexual orientation, age, marital status, veteran status, or disability status.

To make sure our process is as unbiased as we can make it, please ensure that you have removed any personally identifiable information (e.g., your name, website, email, Github username, etc.) from the code. We want to make sure we assess each submission on the quality of its code and only that.

### Must-Haves
These are the requirements your submission must fulfill to be considered correct.

* Send us all source code for test and production, and any artifacts (build scripts, etc.)
* Your code is something you'd be comfortable putting in production and having your team maintaining
* Your code must pass the supplied test harness using the default parameters
* Your code builds on the latest Ubuntu Docker image (if you need a runtime, e.g., Python, JVM or Ruby, let us know which one)
* Your code is tested in some automated fashion at unit and integration levels
* Source control history (e.g., the `.git` directory or a link to Github)

