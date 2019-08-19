var closing = {}
closing.date = ko.observable(false)
closing.dataMaster = ko.observableArray([])
closing.dateEnd = ko.observable()
closing.datalastdate = ko.observableArray([])
closing.debetmaxfive = ko.observableArray([])
closing.creditmaxfive = ko.observableArray([])
closing.lastDate = ko.observable()
closing.TitelFilter = ko.observable(" Hide Filter")
closing.DatePageBar = ko.observable()
closing.swaltext = ko.observable("")
closing.textSearch = ko.observable("")

closing.getDateNow = function () {
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    closing.DatePageBar(page)
}

closing.getData = function (callback) {
    model.Processing(true)
    var url = "/closing/getdatacoatoclosing"
    var datepicker = $('#dateClosing').data('kendoDatePicker').value();
    var now = new Date();
    var dateStart, dateEnd = ""
    if (closing.date() == false) {
        var years = moment(now).format("YYYY")
        var Month = moment(now).format("MM")
        var newDate = years + "-" + Month + "-" + "01"
        dateStart = new Date(newDate)
        dateEnd = new Date()
        closing.dateEnd(dateEnd)
    } else {
        var years = moment(datepicker).format("YYYY")
        var Month = moment(datepicker).format("MM")
        var newDate = years + "-" + Month + "-" + "01"
        dateStart = new Date(newDate)
        if (Month == moment(now).format("MM")) {
            dateEnd = new Date()
            closing.dateEnd(dateEnd)
        } else {
            dateEnd = (new Date(datepicker.getFullYear(), datepicker.getMonth() + 1, 0, 23, 59, 59));
            closing.dateEnd(dateEnd)
        }
    }

    var param = {
        DateStart: dateStart,
        DateEnd: dateEnd,
        Filter: true,
        TextSearch: closing.textSearch().toUpperCase(),
    }
    ajaxPost(url, param, function (res) {
        model.Processing(false)
        if (res.IsError) {
            swal("Search Not Found!", res.Message, "warning")
            $('#textSearch').val("")
            return
        }

        closing.dataMaster(res.Data)
        var arraydebet = _.take(_.sortBy(closing.dataMaster(), function (e) {
            if (e.Debet_Credit == "DEBET") {
                return -e.Transaction;
            }
        }), 5)

        var Listarraydebet = [];
        _.each(arraydebet, function (v, i) {
            Listarraydebet.push({
                code: v.ACC_Code,
                name: v.Account_Name
            });
        });

        closing.debetmaxfive(Listarraydebet)

        var arraykredit = _.take(_.sortBy(closing.dataMaster(), function (e) {
            if (e.Debet_Credit == "CREDIT") {
                return -e.Transaction;
            }
        }), 5)

        var Listarraykredit = [];
        _.each(arraykredit, function (v, i) {
            Listarraykredit.push({
                code: v.ACC_Code,
                name: v.Account_Name
            });
        });

        closing.creditmaxfive(Listarraykredit)

        callback()
    }, function () {
        swal({
            title:"Error!",
            text: "Unknown error, please try again",
            type: "error",
            confirmButtonColor:"#3da09a"})
    })
}

closing.getLastDate = function () {
    var url = "/closing/getlastdate"
    ajaxPost(url, {}, function (res) {
        if (res.IsError === "true") {
            swal({
                title:"Error!", 
                text:res.Message,
                type: "error",
                confirmButtonColor:"#3da09a"})
            return
        }
        if (res.Data.length == 0) {
            closing.lastDate(new Date(01, 0, 01))
        } else {
            var years = moment(res.Data[0].lastdate).year()
            var month = moment(res.Data[0].lastdate).month() + 1
            closing.lastDate(new Date(years, month, 01))
            closing.datalastdate(res.Data)
        }

    })
}

closing.rendergrid = function () {
    var data = closing.dataMaster()
    var columns = [{
        field: '',
        title: 'Date',
        width: 80,
        template: "1 - #= moment(closing.dateEnd()).format('D MMM YYYY')#",
        footerTemplate: "<div style='font-size: 15px;'>TOTAL</div",
    }, {
        field: 'ACC_Code',
        title: 'Acc Code',
        width: 50,

    }, {
        field: 'Account_Name',
        title: 'Account Name',
        width: 200,
    }, {
        field: 'Beginning',
        title: 'Beginning',
        width: 100,
        attributes: {
            style: "text-align:right;"
        },
        template: "#=ChangeToRupiah(Beginning)#",
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
    }, {
        field: 'Transaction',
        title: 'Transaction',
        width: 100,
        attributes: {
            style: "text-align:right;"
        },
        template: "#=ChangeToRupiah(Transaction)#",
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",
    }, {
        field: 'Ending',
        title: 'Ending',
        width: 100,
        attributes: {
            style: "text-align:right;"
        },
        template: "#=ChangeToRupiah(Ending)#",
        footerTemplate: "<div style='text-align:right; font-size: 15px;'>#= ChangeToRupiah(sum) #</div>",

    }]

    $('#gridClosing').html("")
    $('#gridClosing').kendoGrid({
        dataSource: {
            data: data,
            sort: [{
                field: "ACC_Code",
                dir: "asc"
            }],
            aggregate: [{
                field: "Ending",
                aggregate: "sum"
            }, {
                field: "Beginning",
                aggregate: "sum"
            }, {
                field: "Transaction",
                aggregate: "sum"
            }],
        },
        scrollable: true,
        columns: columns
    })
}

closing.refresh = function () {
    closing.date(true)
    closing.init()
}
closing.search = function() {
    // closing.refresh(true)
    closing.getData(function () {
        closing.rendergrid()
    })
 }

closing.closedprocess = function () {
    var url = "/closing/saveandclosing"
    var datepicker = $('#dateClosing').data('kendoDatePicker').value();
    var now = new Date();
    var dateStart, dateEnd = ""
    if (closing.date() == false) {
        var years = moment(now).format("YYYY")
        var Month = moment(now).format("MM")
        var newDate = years + "-" + Month + "-" + "01"
        dateStart = new Date(newDate)
        dateEnd = new Date()
        closing.dateEnd(dateEnd)
    } else {
        var years = moment(datepicker).format("YYYY")
        var Month = moment(datepicker).format("MM")
        var newDate = years + "-" + Month + "-" + "01"
        dateStart = new Date(newDate)
        if (Month == moment(now).format("MM")) {
            dateEnd = new Date()
            closing.dateEnd(dateEnd)
        } else {
            dateEnd = (new Date(datepicker.getFullYear(), datepicker.getMonth() + 1, 0, 23, 59, 59));
            closing.dateEnd(dateEnd)
        }
    }
    var param = {
        DateStart: dateStart,
        DateEnd: dateEnd,
        Data: closing.dataMaster(),
        Periode: "1 -" + moment(closing.dateEnd()).format('D MMM YYYY'),
        SumBegining: $("#gridClosing").data("kendoGrid").dataSource.aggregates().Beginning.sum,
        SumTransaction: $("#gridClosing").data("kendoGrid").dataSource.aggregates().Transaction.sum,
        SumEndig: $("#gridClosing").data("kendoGrid").dataSource.aggregates().Ending.sum
    }

    var textswal = ""
    textswal += "<table width='100%' border='1'>"
    textswal += "<tr><th colspan='2' style='text-align:center'>You will closing at Date 1 - " + moment(closing.dateEnd()).format('D MMM YYYY') + "</th></tr>"
    textswal += " <tr> <td>Debet Rank 5</td> <td>Credit Rank 5</td></tr>"
    textswal += "<tr style='text-align:left;font-size:10px;'><td>"

    _.each(closing.debetmaxfive(), function (v, i) {
        textswal += v.code + " - " + v.name + "<br/>"
    });

    textswal += "</td><td>"

    _.each(closing.creditmaxfive(), function (v, i) {
        textswal += v.code + " - " + v.name + "<br/>"
    });

    textswal += "</td></tr></table>"
    closing.swaltext(textswal)

    swal({
        html: true,
        title: "Are you sure?",
        text: closing.swaltext(),
        type: "warning",
        showCancelButton: true,
        confirmButtonColor: "#3da09a",
        confirmButtonText: "Yes, do it!",
        cancelButtonText: "No!",
        closeOnConfirm: true,
        closeOnCancel: false
    },

    function (isConfirm) {
        model.Processing(true)
        if (isConfirm) {
            ajaxPost(url, param, function (res) {
                model.Processing(false)
                if (res.IsError == true) {
                    return swal({
                        title:"Error!",
                        text: res.Message, 
                        type:"error",
                        confirmButtonColor:"#3da09a"})
                }
                setTimeout(function () {
                    swal({
                        title:"Success!",
                        text: res.Message,
                        type: "success",
                        confirmButtonColor:"#3da09a"});
                    location.reload();
                }, 100)
            })

        } else {
            setTimeout(function () {
                swal({
                    title:"Cancelled", 
                    type:"error",
                    confirmButtonColor:"#3da09a"});
                location.reload();
            }, 100)
        }
    });
}

closing.init = function () {
    // closing.search()
    closing.getLastDate()
    closing.getData(function () {
        closing.rendergrid()
        closing.getDateNow()
    })
}

$(function () {
    closing.init()
})