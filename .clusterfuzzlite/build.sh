#!/bin/bash -eu

go install github.com/AdamKorcz/go-118-fuzz-build@latest
go get github.com/AdamKorcz/go-118-fuzz-build/utils

go get github.com/RICSecLab/mercari-microservices-example/services/authority
go get github.com/RICSecLab/mercari-microservices-example/services/customer
go get github.com/RICSecLab/mercari-microservices-example/services/catalog
go get github.com/RICSecLab/mercari-microservices-example/services/item

compile_native_go_fuzzer github.com/RICSecLab/mercari-microservices-example/services/authority/fuzz FuzzAuthority FuzzAuthority
compile_native_go_fuzzer github.com/RICSecLab/mercari-microservices-example/services/customer/fuzz FuzzCreateCustomer FuzzCreateCustomer
compile_native_go_fuzzer github.com/RICSecLab/mercari-microservices-example/services/catalog/fuzz FuzzCatalog FuzzCatalog
compile_native_go_fuzzer github.com/RICSecLab/mercari-microservices-example/services/item/fuzz FuzzItem FuzzItem

