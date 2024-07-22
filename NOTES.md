This project contains a simple possible implementation of the payment authorization gateway.

The problem was split into the following components (Go modules), in alphabetic order:

1. `bank`  
   this module provides the remote call capability, encapsulating any connections with the authorizing bank  
   `bank` is mockable for the purposes of testing other modules (only the interface `bank.Bank` is exported)

2. `common`  
   this module contains a few project-wide types (like `PaymentRequest`) and utility functions,
   both for the code and for tests

3. `payment`  
   this component provides the "/pay" endpoint; it accepts and validates incoming payment requests,
   performs the remote call to the authorizing bank, via `bank`, and returns a result,
   also storing the latter in the `storage`  
   `payment` is mockable for the purposes of testing other modules (only the interface `payment.Payment` is exported)

4. `storage`  
   this component encapsulates any data-storing/retrieval capabilities, to keep the successful payment (un)authorizations  
   it also provides the "/recall" endpoint used to retrieve earlier payments by id  
   `storage` is mockable for the purposes of testing other modules (only the interface `storage.Storage` is exported)

Additionally, one endpoint is provided by the package `main`, contained in the **redirects.go** file, "/make\_payment",
which is but a wrapper that converts search query parameters into the objects of type `PaymentRequest` and redirects them
to "/pay". It is used mostly for manual testing purposes, in particular, the **html/index.html** sends the results of
the "Make a payment" form to this endpoint.

Other directories of interest are:
* **html**, that contains the file **index.html**; it is a very simple and minimalistic
webpage used for manual tests  
  (NB: on that page, the **amount** field accepts a number in the currency units, not in cents; thus, to send a request
with `"amount":100`, one needs to enter `1.00` in the field, or simply `1`)
* **other_tests**, that contains a simple script **integration.sh** that runs a few scenarios on a real running project
with a live mountebank simmulator.

## Payment processing

**payment/payment.go** contains an implementation of the logic behind the "/pay" endpoint. Few things left noting are:

1. The unmarshaling and validation of the incoming requests is almost automatic, performed by the (almost) standard tools,
`encoding/json` and `validator/v10`, the latter required by the Gin framework.

2. All the communications with the bank and local storage are abstracted out by using packages `bank` and `storage`, respectively.
This would allow for effective mocking inputs and outputs if the package `payment` was tested properly. For simplicity and brevity,
it is currently not, although **integration.sh** effectively runs most of the typical scenarios how `payment` can possibly work,
allowing to observe the `package`'s behaviour.

3. Again for simplicity, most errors produced during any request processing are returned to the customer as is, or with some
minor wrapping. In a real world application we would to need to proces them further on, possibly turning into various sorts of
"request can't be processed at this time" or "something went wrong".

The scripts and all the files in this repo might have their permissions bit `x` set.
It's because this project was developed in an Ubuntu environment but hosted on a Windows machine, on a disk with a Windows file system,
so the permission flags might get crumpled a little bit.

Please let me know in case of any other questions.
