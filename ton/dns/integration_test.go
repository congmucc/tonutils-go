package dns

import (
	"context"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
	"testing"
	"time"
)

var client = liteclient.NewConnectionPool()

var api = func() *ton.APIClient {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := client.AddConnectionsFromConfigUrl(ctx, "https://ton-blockchain.github.io/global.config.json")
	if err != nil {
		panic(err)
	}

	return ton.NewAPIClient(client)
}()

func TestDNSClient_Resolve(t *testing.T) {
	root, err := RootContractAddr(api)
	if err != nil {
		t.Fatal(err)
	}

	cli := NewDNSClient(api, root)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	ctx = client.StickyContext(ctx)

	d, err := cli.Resolve(ctx, "foundation.ton")
	if err != nil {
		t.Fatal(err)
	}

	// wal := d.GetWalletRecord()

	iData, err := d.GetNFTData(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if iData.OwnerAddress.String() != "EQCdqXGvONLwOr3zCNX5FjapflorB6ZsOdcdfLrjsDLt3Fy9" {
		t.Fatal("owner diff", iData.OwnerAddress.String())
	}

	// if wal.String() != "EQA0i8-CdGnF_DhUHHf92R1ONH6sIA9vLZ_WLcCIhfBBXwtG" {
	// 	t.Fatal("wallet record diff", wal.String())
	// }
}

func TestDNSClient_ResolveSub(t *testing.T) {
	root, err := RootContractAddr(api)
	if err != nil {
		t.Fatal(err)
	}

	cli := NewDNSClient(api, root)
	_, err = cli.Resolve(context.Background(), "aa.casino.ton")
	if err != nil {
		if err != ErrNoSuchRecord {
			t.Fatal(err)
		}
	}

	_, err = cli.Resolve(context.Background(), "buchbahpih.ton")
	if err != nil {
		if err != ErrNoSuchRecord {
			t.Fatal(err)
		}
	}
}
