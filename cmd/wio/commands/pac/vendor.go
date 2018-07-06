package pac

type VendorOp int

type Vendor struct {
    Context *cli.Context
    Op VendorOp
}

func (cmd Vendor) GetContext() *cli.Context {
    return cmd.Context;
}

func (cmd Vendor) Execute() error {

}
