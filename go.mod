module github.com/ents-source/member-label-print

go 1.24.0

require (
	github.com/ents-source/go-amember-api v0.0.0-20260523221733-013589c416be
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/kenshaw/escpos v0.0.0-20221114190919-df06b682a8fc
	github.com/tarm/serial v0.0.0-20180830185346-98f6abe2eb07
)

require golang.org/x/sys v0.37.0 // indirect

//replace (
//	github.com/ents-source/go-amember-api => ../go-amember-api
//)
