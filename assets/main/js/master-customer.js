var customer = {}
customer.newRecord = function () {
    return {
        _id: ko.observable(''),
        ID: ko.observable(''),
        Kode: ko.observable(''),
        Name: ko.observable(''),
        Address: ko.observable(''),
        City: ko.observable(''),
        NoTelp: ko.observable(''),
        Owner: ko.observable(''),
        Bank: ko.observable(''),
        AccountNo: ko.observable(''),
        NPWP: ko.observable(''),
        Email: ko.observable(''),
        Type: ko.observable(''),
        TrxCode: ko.observable(''),
        VATReg: ko.observable(''),
        PaymentTerm : ko.observable(''),
        SalesCode: ko.observable(''),
        DepartmentCode : ko.observable(''),
        LastDate : ko.observable(''),
        Limit : ko.observable(0),
    }
}
customer.record = ko.mapping.fromJS(customer.newRecord())
customer.textTitle = ko.observable("")
customer.button = ko.observable("Save")
model.ModeName = ko.observable("#EditCustomer");
customer.dataCustomer = ko.observableArray([]);
customer.idCustomer = ko.observable("");
customer.dataSearch = ko.observable("");
customer.TitelFilter = ko.observable(" Hide Filter");
customer.DatePageBar = ko.observable();
customer.SupplierNumber = ko.observable();
customer.sequenceNumberSupplier = ko.observable();
customer.CustomerNumber = ko.observable();
customer.sequenceNumberCustomer = ko.observable();
customer.URLBack = ko.observable("/transaction/invoice")
customer.TextBack = ko.observable("Link to Invoice")
customer.valCodeBilling = ko.observable();
customer.FlagTab = ko.observable(0);
customer.detailListRecord = {
       _id: ko.observable(''),
        ID: ko.observable(''),
        Kode: ko.observable(''),
        Name: ko.observable(''),
        Address: ko.observable(''),
        City: ko.observable(''),
        NoTelp: ko.observable(''),
        Owner: ko.observable(''),
        Bank: ko.observable(''),
        AccountNo: ko.observable(''),
        NPWP: ko.observable(''),
        Email: ko.observable(''),
        Type: ko.observable(''),
        TrxCode: ko.observable(0),
        VATReg: ko.observable(''),
        PaymentTerm : ko.observable(''),
        SalesCode: ko.observable(''),
        DepartmentCode : ko.observable(''),
        LastDate : ko.observable(''),
        Balance: ko.observable(''),
        Limit: ko.observable(0)
}

customer.dataMasterDepartement = ko.observableArray([])
customer.codeDepartementForDropDown = ko.observableArray([])
customer.dataMasterSales = ko.observableArray([])
customer.codeSalesForDropDown = ko.observableArray([])

customer.getDateNow = function () {
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    customer.DatePageBar(page)
}
customer.getDataCustomer = function (callback) {
    if (customer.FlagTab()==0) {
        var param = {
            TextSearchCode: $("#filterCodeCS").val(),
            TextSearchName: $("#filterNameCS").val(),
        }
    } else {
        var param = {
            TextSearchCode: $("#filterCodeSP").val(),
            TextSearchName: $("#filterNameSP").val(),
        }
    }

    model.Processing(true)
    ajaxPost('/master/getdatacustomer', param, function (res) {
        if (res.Total === 0) {
            swal({
                title:"Error!",
                text: res.Message, 
                type:"error",
                confirmButtonColor:"#3da09a"
            })
            return
        }

        customer.dataCustomer(res.Data)
        model.Processing(false)
        callback()
    }, function () {
        swal({
            title:"Error!",
            text: "Unknown error, please try again",
            type: "error",
            confirmButtonColor:"#3da09a"
        })
    })
}

customer.dataCustomer = ko.observableArray([])
customer.getEditCustomer = function (callback) {

    ajaxPost('/master/geteditcustomer', {}, function (res) {
        if (res.Total === 0) {
            swal({
                title:"Error!",
                text: res.Message, 
                type:"error",
                confirmButtonColor:"#3da09a"
        })
            return
        }

        customer.dataCustomer(res.Data)

        callback()
    }, function () {
        swal({
            title:"Error!",
            text: "Unknown error, please try again",
            type: "error",
            confirmButtonColor:"#3da09a"
        })
    })
}

// === Show hidden field when exporting
var exportExcelFlag = 0;
// ====================================
customer.renderGrid = function (search) {
    var data = _.filter(customer.dataCustomer(), function (item) {
        return item.Type == "CUSTOMER"
    });
    if (search != undefined) {
        var results = _.filter(data, function (item) {
            return item.AccountNo.indexOf(search) > -1 || item.Address.indexOf(search) > -1 || item.Bank.indexOf(search) > -1 || item.City.indexOf(search) > -1 || item.Email.indexOf(search) > -1 || item.Kode.indexOf(search) > -1 || item.NPWP.indexOf(search) > -1 || item.Name.indexOf(search) > -1 || item.NoTelp.indexOf(search) > -1 || item.Owner.indexOf(search) > -1 || item.Type.indexOf(search) > -1 || item._id.indexOf(search) > -1
        });
        data = results
    }
    if (search == "") {
        data
    }

    if (typeof $('#gridcustomerlist').data('kendoGrid') !== 'undefined') {
        $('#gridcustomerlist').data('kendoGrid').setDataSource(new kendo.data.DataSource({
            data: data,
            pageSize: 25
        }))
        return
    }
    console.log(data)
    var columns = [
    {
    field: 'Kode',
    title: 'Code',
    width: 150,
    template: function (d) {
        return '<a class="onclickJurnal" href="javascript:void(0)" onclick="customer.detailList(\'' + d._id + '\')">' + d.Kode + '</a>'
        },
    }, {
        field: 'Name',
        title: 'Company Name'
    }, {
        field: 'Address',
        title: 'Address',
        width: 150,
        hidden: true
    }, {
        field: 'City',
        title: 'City',
        width: 150,
        hidden: true
    }, {
        field: 'NoTelp',
        title: 'Phone',
        width: 150,
        hidden: true
    }, {
        field: 'Owner',
        title: 'Owner',
        width: 150,
        hidden: true
    }, {
        field: 'Bank',
        title: 'Bank',
        width: 150,
        hidden: true
    }, {
        field: 'AccountNo',
        title: 'Bank Account Number',
        width: 150,
        hidden: true
    }, {
        field: 'VATReg',
        title: 'VAT Reg. No.',
        width: 150,
        hidden: true
    }, {
        field: 'PaymentTerm',
        title: 'Payment Term',
        width: 150,
        hidden: true
    }, {
        field: 'Email',
        title: 'Email',
        width: 150,
        hidden: true
    }, {
        field: 'SalesCode',
        title: 'Sales Code',
        width: 150,
        hidden: true
    }, {
        field: 'DepartmentCode',
        title: 'Department Code',
        width: 150,
        hidden: true
    }, {
        field: 'Limit',
        title: 'Limit',
        width: 150,
        hidden: true
    },{
        title: 'Action',
        width: 110,
        template: "<a href=\"javascript:customer.detailBillingCS('#: Kode #')\"  data-backdrop=\"static\" style=\"width:12px;\" class=\"btn btn-xs btn-primary\"><i class=\"fa fa-info\"></i></a>    <a href=\"javascript:customer.editCustomer('#: _id #')\" data-target=\".EditCustomer\" data-backdrop=\"static\" class=\"btn btn-xs btn-warning\"><i class=\"fa fa-pencil\"></i></a>&nbsp;" + "#if(1==1){#" + "<a href=\"javascript:customer.delete('#: _id #', '#: Kode #')\" class=\"btn btn-xs btn-danger\"><span class='glyphicon glyphicon-trash'></span>#}#"
    }]

    $('#gridcustomerlist').kendoGrid({
        dataSource: {
            data: data,
            pageSize: 10,
        },
        pageable: {
            // refresh: true,
            pageSizes: true,
            buttonCount: 5
        },
        sortable: true,
        height: 400,
        width: 140,
        filterable: false,
        scrollable: true,
        columns: columns,
        excelExport: function(e) {
            if (exportExcelFlag == 0)
            {
                for (var i=2; i<=12; i++)
                    e.sender.showColumn(i);
                setTimeout(() => {
                    exportExcelFlag = 1;
                    e.sender.saveAsExcel();
                });
                e.preventDefault();
            } else if (exportExcelFlag == 1) {
                ProActive.kendoExcelRender(e, "Customer", function(row, sheet){
    
                });
                setTimeout(() => {
                    exportExcelFlag = 2;
                    e.sender.saveAsExcel();
                });
            } else if (exportExcelFlag == 2) {
                for (var i=2; i<=12; i++)
                    e.sender.hideColumn(i);
                customer.renderGrid();
                exportExcelFlag = 0;
                e.preventDefault();
            }
        },
    })
}
customer.renderGridSupplier = function (search) {
    var data = _.filter(customer.dataCustomer(), function (item) {
        return item.Type == "SUPPLIER";
    });

    if (search != undefined) {
        var results = _.filter(data, function (item) {
            return item.AccountNo.indexOf(search) > -1 || item.Address.indexOf(search) > -1 || item.Bank.indexOf(search) > -1 || item.City.indexOf(search) > -1 || item.Email.indexOf(search) > -1 || item.Kode.indexOf(search) > -1 || item.NPWP.indexOf(search) > -1 || item.Name.indexOf(search) > -1 || item.NoTelp.indexOf(search) > -1 || item.Owner.indexOf(search) > -1 || item.Type.indexOf(search) > -1 || item._id.indexOf(search) > -1
        });
        data = results
    }
    if (search == "") {
        data
    }

    if (typeof $('#gridsupplier').data('kendoGrid') !== 'undefined') {
        $('#gridsupplier').data('kendoGrid').setDataSource(new kendo.data.DataSource({
            data: data,
            pageSize: 25
        }))
        return
    }

    var columns = [
        {
        field: 'Kode',
        title: 'Code',
        width: 150,
        template: function (d) {
            return '<a class="onclickJurnal" href="javascript:void(0)" onclick="customer.detailList(\'' + d._id + '\')">' + d.Kode + '</a>'
            },
        }, {
            field: 'Name',
            title: 'Company Name'
        }, {
            field: 'Address',
            title: 'Address',
            width: 150,
            hidden: true
        }, {
            field: 'City',
            title: 'City',
            width: 150,
            hidden: true
        }, {
            field: 'NoTelp',
            title: 'Phone',
            width: 150,
            hidden: true
        }, {
            field: 'Owner',
            title: 'Owner',
            width: 150,
            hidden: true
        }, {
            field: 'Bank',
            title: 'Bank',
            width: 150,
            hidden: true
        }, {
            field: 'AccountNo',
            title: 'Bank Account Number',
            width: 150,
            hidden: true
        }, {
            field: 'VATReg',
            title: 'VAT Reg. No.',
            width: 150,
            hidden: true
        }, {
            field: 'PaymentTerm',
            title: 'Payment Term',
            width: 150,
            hidden: true
        }, {
            field: 'Email',
            title: 'Email',
            width: 150,
            hidden: true
        }, {
            field: 'SalesCode',
            title: 'Sales Code',
            width: 150,
            hidden: true
        }, {
            field: 'DepartmentCode',
            title: 'Department Code',
            width: 150,
            hidden: true
        },{
            title: 'Action',
            width: 110,
            template: "<a href=\"javascript:customer.detailBillingSPNonInv('#: Kode #')\"  data-backdrop=\"static\" style=\"width:12px;\" class=\"btn btn-xs btn-primary\"><i class=\"fa fa-info\"></i></a>    <a href=\"javascript:customer.editCustomer('#: _id #')\" data-target=\".EditCustomer\" data-backdrop=\"static\" class=\"btn btn-xs btn-warning\"><i class=\"fa fa-pencil\"></i></a>&nbsp;" + "#if(1==1){#" + "<a href=\"javascript:customer.delete('#: _id #', '#: Kode #')\" class=\"btn btn-xs btn-danger\"><span class='glyphicon glyphicon-trash'></span>#}#"
        }]
        

    $('#gridsupplier').kendoGrid({
        dataSource: {
            data: data,
            pageSize: 10,

        },
        pageable: {
            // refresh: true,
            pageSizes: true,
            buttonCount: 5
        },
        sortable: true,
        height: 400,
        width: 140,
        filterable: false,
        scrollable: true,
        columns: columns,
        excelExport: function(e) {
            if (exportExcelFlag == 0)
            {
                for (var i=2; i<=12; i++)
                    e.sender.showColumn(i);
                setTimeout(() => {
                    exportExcelFlag = 1;
                    e.sender.saveAsExcel();
                });
                e.preventDefault();
            } else if (exportExcelFlag == 1) {
                ProActive.kendoExcelRender(e, "Supplier", function(row, sheet){
    
                });
                setTimeout(() => {
                    exportExcelFlag = 2;
                    e.sender.saveAsExcel();
                });
            } else if (exportExcelFlag == 2) {
                for (var i=2; i<=12; i++)
                    e.sender.hideColumn(i);
                customer.renderGrid();
                exportExcelFlag = 0;
                e.preventDefault();
            }
        },
    })
}



customer.exportExcel = function () {
    if (choosen == "SUP")
        $("#gridsupplier").getKendoGrid().saveAsExcel();
    else
        $("#gridcustomerlist").getKendoGrid().saveAsExcel();
}

customer.detailList = function(id) {
    //console.log(id);
    var data = _.filter(customer.dataCustomer(), function (item) {
        return item._id == id
    });
    console.log(data)
    var d = customer.detailListRecord
    d.ID(data[0]._id)
    d.Kode(data[0].Kode)
    d.Name(data[0].Name)
    d.Address(data[0].Address)
    d.City(data[0].City)
    d.NoTelp(data[0].NoTelp)
    d.Owner(data[0].Owner)
    d.Bank(data[0].Bank)
    d.AccountNo(data[0].AccountNo)
    d.NPWP(data[0].NPWP)
    d.Email(data[0].Email)
    d.Type(data[0].Type)
    d.TrxCode(data[0].TrxCode)
    d.VATReg(data[0].VATReg)
    d.PaymentTerm(data[0].PaymentTerm)
    if (data[0].SalesCode == undefined) {
        d.SalesCode("N/A")
    } else {
        d.SalesCode(data[0].SalesCode)
    }
    if (data[0].DepartementCode == "") {
        d.DepartmentCode("N/A")
    } else {
        d.DepartmentCode(data[0].DepartementCode)
    }
    
    d.LastDate(data[0].LastDate)
    d.Limit(ChangeToRupiah(data[0].Limit))

    customer.getBalance(data[0].Kode, data[0].Type);
}

customer.getBalance = function(code, type) {
    var p = {
        CustomerCode : code
    }

    ajaxPost((type == "CUSTOMER" ? '/master/GetCustomerBalance' : '/master/GetSupplierBalance').toLowerCase(), p, function (res) {
        if (!res.IsError) {
            customer.detailListRecord.Balance(ChangeToRupiah(res.Data.BalanceIDR))
        }
    })
}

customer.clear = function(){
    var d = customer.detailListRecord
    d.ID("")
    d.Kode("")
    d.Name("")
    d.Address("")
    d.City("")
    d.NoTelp("")
    d.Owner("")
    d.Bank("")
    d.AccountNo("")
    d.NPWP("")
    d.Email("")
    d.Type("")
    d.TrxCode("")
    d.VATReg("")
    d.PaymentTerm("")
    d.SalesCode("")
    d.DepartmentCode("")
    d.LastDate("")
    d.Balance("")
    d.Limit("")
}

customer.getDataSales = function () {
    model.Processing(true)

    ajaxPost('/master/getdatasales', {}, function (res) {

        if (res.Total === 0) {
            swal("Error!", res.Message, "error")
            return
        }
        customer.dataMasterSales(res.Data)
        var DataAccount = res.Data
        for (i in DataAccount) {
            DataAccount[i].SalesCode = DataAccount[i].SalesID
            DataAccount[i].SalesName = DataAccount[i].SalesID + " - "+ DataAccount[i].SalesName
        }
        customer.codeSalesForDropDown(DataAccount)
        model.Processing(false)
    })
}

customer.getDataDepartment = function () {
    model.Processing(true)

    ajaxPost('/transaction/getdatadepartment', {}, function (res) {

        if (res.Total === 0) {
            swal("Error!", res.Message, "error")
            return
        }
        customer.dataMasterDepartement(res.Data)
        // console.log(customer.dataMasterDepartement())
        var DataAccount = res.Data
        for (i in DataAccount) {
            DataAccount[i].DepartmentCode = DataAccount[i].DepartmentCode
            DataAccount[i].DepartmentName = DataAccount[i].DepartmentCode + " - "+ DataAccount[i].DepartmentName
        }
        customer.codeDepartementForDropDown(DataAccount)
        //console.log(customer.codeDepartementForDropDown())

        model.Processing(false)
    })
}

customer.addNew = function () {
    customer.button("Save")
    if (customer.record.Type() == "CUSTOMER") {
        ko.mapping.fromJS(customer.newRecord(), customer.record)
        customer.record.Type("CUSTOMER")
        $(".custhide").show()
    } else {
        ko.mapping.fromJS(customer.newRecord(), customer.record)
        customer.record.Type("SUPPLIER")
        $(".custhide").hide()
    }
    customer.record.Kode("")
}

function FilterAlphanum(str, charset) {
    if (!str) return str;
    charset = charset || "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
    var re = "";
    for(var i=0; i < str.length; i++) {
        if (charset.includes(str.substr(i,1))) {    
            re += str.substr(i,1);
        }
    }
    return re;
}

customer.saveData = function () {
    customer.record.ID(customer.record._id())
    var param = ko.mapping.toJS(customer.record)
    param.Name = param.Name.toUpperCase();
    param.TrxCode = parseInt(param.TrxCode)

    myLimit = param.Limit+""
    param.Limit = parseFloat(myLimit.replace(/\,/g,''))   

    var struct = {
        Data: param,
    }
    var url = "/master/insertnewcustomer"

    //validation
    var vldName = FilterAlphanum(param.Name);
    if (vldName == "") {
        return swal({
            title:'Warning!', 
            text:"You haven't fill the Company Name", 
            type:"info",
            confirmButtonColor:"#3da09a"
        })
    }else if (vldName.length < 3) {
        return swal({
            title:'Warning!', 
            text:"The Company Name is too short (min 3 letters)", 
            type:"info", 
            confirmButtonColor:"#3da09a"
        })
    } else {
    swal({
        title: "Are you sure?",
        text: "You will submit this " + param.Type.toLowerCase(),
        type: "warning",
        showCancelButton: true,
        confirmButtonColor: '#3da09a',
        confirmButtonText: "Yes, do it!",
        cancelButtonText: "No!",
        closeOnConfirm: true,
        closeOnCancel: false
    }, function (isConfirm) {
        if (isConfirm) {
            model.Processing(true)
            ajaxPost(url, struct, function (e) {
                customer.reset()
                model.Processing(false)
                if (e.Status == "OK") {
                    swal({
                        title:"Success!", 
                        text:e.Message, 
                        type:"success",
                        confirmButtonColor:"#3da09a"
                    });
                    customer.init()
                    $('#AddNewCustomer').modal('hide');
                    if (customer.record._id() != "") {
                        setTimeout(() => {
                            customer.detailList(customer.record._id())
                        }, 100);
                    }
                } else {
                    swal({
                        title:"Error!",
                        text: e.Message,
                        type: "error",
                        confirmButtonColor:"#3da09a"
                    });
                }
            })
        } else {
            swal({
                title: "Cancelled", 
                type: "error",
                closeOnConfirm: true,
                confirmButtonColor:'#3da09a',
                confirmButtonText: 'OK'});
        }
    });
    }
}

customer.reset = function () {
    kodepelanggan = $("#kodepelanggan").val("")
    name = $("#name").val("")
    address = $("#address").val("")
    city = $("#city").val("")
    notelp = $("#telp").val("")
    owner = $("#owner").val("")
    bank = $("#bank").val("")
    accountno = $("#account").val("")
    npwp = $("#npwp").val("")
    email = $("#email").val("")
}

customer.editCustomer = function (id) {
    var data = _.find(customer.dataCustomer(), function (item) {
        return item._id == id;
    });   
    ko.mapping.fromJS(data, customer.record)
    $("#AddNewCustomer").modal("show")
    customer.button("Update")
}

customer.detailBillingCS = function (code) {
    $("#inventoryradio").hide()
    $("#titleBilling").text("Customer Detail Billing : "+ code)
    customer.valCodeBilling(code)
    var columns = [{
        field: 'DateStr',
        title: 'Date',
        width: 80,
        format: "{0:d-MMM-yyyy}",
    }, {
        field: 'Type',
        title: 'Type Billing',
        width: 100,
    }, {
        field: 'DocumentNumber',
        title: 'Document Number',
        width: 150,
    }, {
        title: 'Total (IDR)',
        width: 100,
        template : function(d){
            if (d.Total == 0) return "--"
            else return ChangeToRupiah(d.Total)
        }
    },{
        title: 'Paid (IDR)',
        width: 100,
        template : function(d){
            if (d.Paid == 0) return "--"
         else return ChangeToRupiah(d.Paid)
        }
    },{
        title: 'Balance',
        field: 'Balance',
        width: 100,
        template : function(d){
            return ChangeToRupiah(d.Balance)
        }
    }]
    $("#DetailBilling").modal("show")
    var param = {
        CustomerCode: code,
    }
    model.Processing(true)
    ajaxPost('/master/getdetailbillingcustomer', param, function (res) {
        Data = _.sortBy(res.Data, function(o){return o.DateStr})

        CustomerBalance = 0;
        for (i = 0; i < Data.length; i++) {
            if (Data[i].Type == "Sales Credit Memo"){
                Data[i].Total = Data[i].Total * -1
            }

            CustomerBalance += Data[i].Total - Data[i].Paid
            Data[i].Balance = CustomerBalance
        }
        model.Processing(false)
        // console.log(Data.reverse())
        $('#griddetailbilling').kendoGrid({
            dataSource: {
                data: Data,
                sort: {
                    field: "DateStr",
                    dir: "asc"
                },
                schema: {
                    model : {
                        fields: {
                            DateStr: {
                                type: "date",
                                parse: function(e) {
                                    return new Date(e)
                                }
                            }
                        }
                    }
                }
            },
            sortable: true,
            height: 400,
            width: 160,
            filterable: false,
            scrollable: true,
            columns: columns
        })
    })
}

customer.detailBillingSPNonInv = function (code) {
    $("#titleBilling").text("Supplier Detail Billing : "+ code)
    customer.valCodeBilling(code)
    $("#inventoryradio").hide()

    var columns = [{
        field: 'DateStr',
        title: 'Date',
        width: 80,
    }, {
        field: 'Type',
        title: 'Type Billing',
        width: 100,
    }, {
        field: 'DocumentNumber',
        title: 'Document Number',
        width: 150,
    }, {
        title: 'Total (IDR)',
        width: 100,
        template : function(d){
            if (d.Total == 0) return "--"
            else return ChangeToRupiah(d.Total)
        }
    },{
        title: 'Paid (IDR)',
        width: 100,
        template : function(d){
            if (d.Paid == 0) return "--"
         else return ChangeToRupiah(d.Paid)
        }
    },{
        title: 'Balance',
        field: 'Balance',
        width: 100,
        template : function(d){
            if (d.Type == "Purchase Invoice (Non Inventory)" || d.Type == "Purchase Invoice") {
                return "<i>" +ChangeToRupiah(d.Balance)+ " *</i>"
            } else {
                return ChangeToRupiah(d.Balance)
            }
        }
    }]
    $("#DetailBilling").modal("show")
    var param = {
        SupplierCode: code,
    }
    model.Processing(true)
    ajaxPost('/master/getdetailbillingsuppliernoninventory', param, function (res) {
        Data = _.sortBy(res.Data, function(o){return o.DateCreated})
        // Data.reverse()
        SupplierBalance = 0;
        for (i = 0; i < Data.length; i++) {
            if (Data[i].Type == "Purchase Credit Memo"){
                Data[i].Total = Data[i].Total * -1
            }
            if (Data[i].Type != "Purchase Invoice (Non Inventory)" && Data[i].Type != "Purchase Invoice") {
                SupplierBalance += Data[i].Total - Data[i].Paid
            }
            Data[i].Balance = SupplierBalance
        }

        model.Processing(false)
        // console.log(Data.reverse())
        $('#griddetailbilling').kendoGrid({
            dataSource: {
                data: Data,
                sort: {
                    field: "DateStr",
                    dir: "asc"
                },
                schema: {
                    model : {
                        fields: {
                            DateStr: {
                                type: "date",
                                parse: function(e) {
                                    return new Date(e)
                                }
                            }
                        }
                    }
                }  
            },
            sortable: true,
            height: 400,
            width: 160,
            filterable: false,
            scrollable: true,
            columns: columns
        })
    })
}

customer.detailBillingSPInv = function (code) {
    $("#titleBilling").text("Supplier Detail Billing : "+ code)
    customer.valCodeBilling(code)
    $("#inventoryradio").show()

    var columns = [{
        field: 'DateStr',
        title: 'Date',
        width: 80,
    }, {
        field: 'Type',
        title: 'Type Billing',
        width: 100,
    }, {
        field: 'DocumentNumber',
        title: 'Document Number',
        width: 150,
    }, {
        title: 'Total (IDR)',
        width: 100,
        template : function(d){
            if (d.Total == 0) return "--"
            else return ChangeToRupiah(d.Total)
        }
    },{
        title: 'Paid (IDR)',
        width: 100,
        template : function(d){
            if (d.Paid == 0) return "--"
         else return ChangeToRupiah(d.Paid)
        }
    },{
        title: 'Balance',
        field: 'Balance',
        width: 100,
        template : function(d){
            return ChangeToRupiah(d.Balance)
        }
    }]
    $("#DetailBilling").modal("show")
    var param = {
        SupplierCode: code,
    }
    model.Processing(true)
    ajaxPost('/master/getdetailbillingsupplierinventory', param, function (res) {
        Data = _.sortBy(res.Data, function(o){return o.CreatedDate})
        Data.reverse()
        SupplierBalance = 0;
        for (i = 0; i < Data.length; i++) {
            if (Data[i].Type == "Purchase Credit Memo"){
                Data[i].Total = Data[i].Total * -1
            }

            SupplierBalance += Data[i].Total - Data[i].Paid
            Data[i].Balance = SupplierBalance
        }

        model.Processing(false)
        // console.log(Data.reverse())
        $('#griddetailbilling').kendoGrid({
            dataSource: {
                data: Data,
            },
            sortable: true,
            height: 400,
            width: 160,
            filterable: false,
            scrollable: true,
            columns: columns
        })
    })
}

customer.checkradio = function (val) {
    if (val == "noninventory") {
        customer.detailBillingSPNonInv(customer.valCodeBilling())
    } else if(val == "inventory") {
        customer.detailBillingSPInv(customer.valCodeBilling())
    }
}

customer.delete = function (id, KodePelanggan) {
    model.Processing(true);
    swal({
        title: "Are you sure to delete " + KodePelanggan + "?",
        text: "Your will not be able to recover this data",
        type: "warning",
        showCancelButton: true,
        confirmButtonClass: "btn-danger",
        confirmButtonText: "Yes, delete it!",
        closeOnConfirm: false
    }, function (res) {
        if (res) {
            var url = "/master/delete";
            var param = {
                id: id
            };
            ajaxPost(url, param, function (data) {
                if (data != "") {
                    swal({
                        title:'Warning',
                        text: data,
                        type: 'error',
                        confirmButtonColor:"#3da09a"
                    });
                    model.Processing(false);
                } else {
                    swal({
                        title:'Success',
                        text: 'User has been deleted!',
                        type: 'success',
                        confirmButtonColor:"3da09a"});
                    customer.init()
                    model.Processing(false);
                }

            }, undefined);
        } else {
            model.Processing(false);
        }
    });
}

customer.tabCustomer = function () {
    customer.FlagTab(0)
    customer.resetFilter()
    customer.textTitle("Customer")
    customer.record.Type("CUSTOMER")
    customer.init()
    choosen = "CUS"
    customer.URLBack("/transaction/invoice")
    customer.TextBack("Link to Invoice")
    $('#textSearch').val("")
    $("#custhide").show()
    customer.clear()
}

customer.tabSupplier = function () {
    customer.FlagTab(1)
    customer.resetFilter()
    customer.textTitle("Supplier")
    customer.record.Type("SUPPLIER")
    customer.init()
    choosen = "SUP"
    customer.URLBack("/transaction/purchaseorder")
    customer.TextBack("Link to PO")
    $('#textSearch').val("")
    customer.dataSearch("")
    $("#custhide").hide()
    customer.clear()
}

customer.resetFilter = function () {
    $("#filterCodeCS").val("")
    $("#filterNameCS").val("")
    $("#filterCodeSP").val("")
    $("#filterNameSP").val("")
}

customer.back = function(){
    location.href=customer.URLBack()
}

customer.search = function () {
    var search = customer.dataSearch()
    if (customer.textTitle() == "Customer") {
        customer.renderGrid(search)
    } else {
        customer.renderGridSupplier(search)
    }
}
customer.maskingMoney = function() {
    $('.currency').inputmask("numeric", {
        radixPoint: ".",
        groupSeparator: ",",
        digits: 2,
        autoGroup: true,
        // prefix: '$', //No Space, this will truncate the first character
        rightAlign: false,
        // oncleared: function () { self.Value(''); }
    });
}
customer.init = function () {
    customer.getDataCustomer(function () {
        customer.renderGrid()
        customer.renderGridSupplier()
        customer.getDateNow()
        customer.getDataDepartment()
        customer.getDataSales()
        customer.maskingMoney
    })
    customer.maskingMoney();
}

$(function () {
    customer.textTitle("Customer")
    customer.record.Type("CUSTOMER")
    customer.init()
    $("#vatreg").mask("99.999.999.9-999.999");
    $('.valnum').keypress(validateNumber);
})