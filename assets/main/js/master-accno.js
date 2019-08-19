var master = {}

master.dataMasterAcc = ko.observableArray([])
master.AccountSelected = ko.observable("")
master.DatePageBar = ko.observable()
master.TitelFilter = ko.observable(" Show Filter");

master.getDateNow = function(){
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    master.DatePageBar(page)
}
master.AccounIdEdited = ko.observable("")
master.AccountNameEdited = ko.observable("")
master.AccountCodeEdited = ko.observable("")
master.AccountMainCodeEdited = ko.observable("")
master.getDataMasterAcc = function (callback) {
    model.Processing(true)
    ajaxPost('/master/getdatacoa', {}, function (res) {
        if (res.Total === 0) {
            swal("Error!", res.Message, "error")
            return
        }
        model.Processing(false)
        master.dataMasterAcc(_.sortBy(res.Data, [function(o){return o.ACC_Code}]))

        callback()
    }, function () {
        swal("Error!", "Unknown error, please try again", "error")
    })
}

master.renderGrid = function () {
    var data = master.dataMasterAcc();
    if (typeof $('.grid-accno').data('kendoGrid') !== 'undefined') {
        $('.grid-accno').data('kendoGrid').setDataSource(new kendo.data.DataSource({
            data: data,
            // pageSize: 25
        }))
        return
    }

    var columns = [
        {
            title: 'ACTION',
            width: 60,
            template: "<button onclick='master.editCOA(\"#: ID #\")' class='btn btn-sm btn-success btn-flat'><i class='fa fa-pencil'></i></button>",
        },{ 
            field: 'ACC_Code', 
            title: 'ACCOUNT NUMBER', 
            width: 140 
        },{ 
            field: 'Account_Name', 
            title: 'ACCOUNT NAME', 
            width: 300 
        },{ 
            field: 'Category', 
            title: 'CATEGORY', 
            width: 200 
        }, { 
            field: 'Debet_Credit', 
            title: 'DEBET/CREDIT', 
            width: 200 
        }
    ]

    $('.grid-accno').kendoGrid({
        dataSource: {
            data: data,
        },
        sortable: true,
        height : 500,
        scrollable:true,
        columns: columns,
        dataBound: function(e) {
          dataView = this.dataSource.view();
            for (var i = 0; i < dataView.length; i++) {
                if (dataView[i].Main_Acc_Code == 0) {
                    var uid = dataView[i].uid;
                    $(".grid-accno").find("tr[data-uid=" + uid + "]").addClass("yellowBackground");
                }
            }
        }

    })
}
master.editCOA = function(id){
    master.AccountNameEdited("")
    master.AccountCodeEdited("")
    master.AccountMainCodeEdited("")
    master.AccounIdEdited("")
    var data = _.find(master.dataMasterAcc(), function(e){ return e.ID == id})
    master.AccountNameEdited(data.Account_Name)
    master.AccountCodeEdited(data.ACC_Code)
    master.AccountMainCodeEdited(data.Main_Acc_Code)
    master.AccounIdEdited(id)
    $('#EditNewModal').modal('show');
}
master.SaveDataCOA = function(){
    if (master.AccountNameEdited()==""){
        return swal("Info!", "Account name is not filled", "warning")
    }
    model.Processing(true)
    var url = "/master/updatecoa"
    var param = {
        Id : master.AccounIdEdited(),
        AccName : master.AccountNameEdited()
    }
    ajaxPost(url, param, function(res){
        model.Processing(false)
        $('#EditNewModal').modal('hide');
        if (res.IsError == true) {
            return swal({
                title :"Error!",
                text : res.Message,
                type : "error",
                confirmButtonColor: '#3da09a'
            })
        }
        master.getDataMasterAcc(function () {
            master.renderGrid()
        })
    })
}
master.checkedData = function () {
    return $('.checkboxgrid').get().filter(function (d) {
        return d.checked
    }).map(function (d) {
        return $(d).attr('data-id')
    })
}
master.DropdownAccount = function(){
    var data = master.dataMasterAcc();
    for (i in data) {
        data[i].ACC_Code = data[i].ACC_Code + ""
        data[i].Account_Name = data[i].ACC_Code +" - " +data[i].Account_Name
    }
    var DataAcc = data
    $('#account-number').width(200).kendoDropDownList({
        dataSource: DataAcc,
        filter: "startswith",
        dataTextField: "Account_Name",
        dataValueField: "ACC_Code",
        optionLabel: "Select Account Number",
        animation: {
            close: {
            effects: "zoom:out",
            duration: 300
            }
        },
        change: function(e){
            var key = this.dataItem();
            master.AccountSelected(key.ACC_Code)
        }
    });
}
master.Refresh = function(){
    var AccSelected = parseInt(master.AccountSelected())
    var Data =  _.find(master.dataMasterAcc(), function (o) { return o.ACC_Code=AccSelected; });
    master.dataMasterAcc([Data])
    master.renderGrid()
}
master.Clear = function(){
    master.init()
    var dropdownlist = $("#account-number").data("kendoDropDownList");
    dropdownlist.value("");
}

master.init = function(){
    master.getDataMasterAcc(function () {
        master.renderGrid()
        master.DropdownAccount()
        master.getDateNow()
    })
}
$(function () {
    master.init()
})
