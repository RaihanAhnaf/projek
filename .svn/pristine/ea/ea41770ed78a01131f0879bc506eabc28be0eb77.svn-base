model.Processing(false)
var glusdmandiri = {}
glusdmandiri.koFilterIsActive = ko.observable(false)
glusdmandiri.koDataPettyCash = ko.observableArray([])
glusdmandiri.refreshGrid = function () {
    glusdmandiri.koFilterIsActive(true)
    glusdmandiri.getDataglusdmandiri(function () {
        glusdmandiri.renderGrid()
    })
}
glusdmandiri.getDataglusdmandiri = function (callback) {
    model.Processing(true)
    var url = "/report/getdataglusdmandiri"
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value();
    var param = {
        DateStart: dateStart,
        DateEnd: dateEnd,
        Filter: glusdmandiri.koFilterIsActive(),
    }
    ajaxPost(url, param, function (res) {
        glusdmandiri.koDataPettyCash(res.Data)
        model.Processing(false)
        callback()
    })
}
glusdmandiri.renderGrid = function () {
    var data = glusdmandiri.koDataPettyCash()
    var columns = [{
        field: 'PostingDate',
        title: 'Posting Date',
        template: function (e) {
            return moment.utc(e.PostingDate).format("DD MMM YYYY")
        },
        width: 70,
    }, {
        field: 'Acc_Code',
        title: 'Acc Code',
        width: 50
    }, {
        field: 'Acc_Name',
        title: 'Account Name',
        width: 100
    }, {
        field: 'Description',
        title: 'Description',
        width: 200,
        footerTemplate: "Total :"
    }, {
        field: 'Debet',
        title: 'Debit',
        width: 100,
        template: "#=ChangeToRupiah(Debet)#",
        footerTemplate: "#= ChangeToRupiah(sum) #"
    }, {
        field: 'Credit',
        title: 'Credit',
        width: 100,
        template: "#=ChangeToRupiah(Credit)#",
        footerTemplate: "#= ChangeToRupiah(sum) #"
    }, {
        field: 'Saldo',
        title: 'Saldo',
        width: 100,
        template: "#=ChangeToRupiah(Credit)#",
        footerTemplate: function (e) {
            return glusdmandiri.saldoCalculation()
        }
    }, ]
    $('#gridGLUSDMandiri').kendoGrid({
        dataSource: {
            data: data,
            aggregate: [{
                field: "Debet",
                aggregate: "sum"
            }, {
                field: "Credit",
                aggregate: "sum"
            }, {
                field: "Saldo",
                aggregate: "sum"
            }],
        },
        height: 500,
        filterable: true,
        scrollable: true,
        columns: columns
    })
}
glusdmandiri.saldoCalculation = function () {
    var data = $('#gridGLUSDMandiri').data('kendoGrid').dataSource.options.data
    var sumDebit = _.sumBy(data, 'Debet')
    var sumCredt = _.sumBy(data, 'Credit')
    var Saldo = sumDebit - sumCredt
    return kendo.toString(Saldo, 'n')
}
glusdmandiri.init = function () {
    glusdmandiri.getDataglusdmandiri(function () {
        glusdmandiri.renderGrid()
    })
}
$(document).ready(function () {
    glusdmandiri.init()
})