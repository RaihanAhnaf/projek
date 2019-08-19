model.Processing(false)
var glidrmandiri = {}
glidrmandiri.koFilterIsActive = ko.observable(false)
glidrmandiri.koDataPettyCash = ko.observableArray([])
glidrmandiri.refreshGrid = function () {
    glidrmandiri.koFilterIsActive(true)
    glidrmandiri.getDataglidrmandiri(function () {
        glidrmandiri.renderGrid()
    })
}
glidrmandiri.getDataglidrmandiri = function (callback) {
    model.Processing(true)
    var url = "/report/getdataglidrmandiri"
    var dateStart = $('#dateStart').data('kendoDatePicker').value();
    var dateEnd = $('#dateEnd').data('kendoDatePicker').value();
    var param = {
        DateStart: dateStart,
        DateEnd: dateEnd,
        Filter: glidrmandiri.koFilterIsActive(),
    }
    ajaxPost(url, param, function (res) {
        glidrmandiri.koDataPettyCash(res.Data)
        model.Processing(false)
        callback()
    })
}
glidrmandiri.renderGrid = function () {
    var data = glidrmandiri.koDataPettyCash()
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
            return glidrmandiri.saldoCalculation()
        }
    }, ]
    $('#gridGLIDRMandiri').kendoGrid({
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
glidrmandiri.saldoCalculation = function () {
    var data = $('#gridGLIDRMandiri').data('kendoGrid').dataSource.options.data
    var sumDebit = _.sumBy(data, 'Debet')
    var sumCredt = _.sumBy(data, 'Credit')
    var Saldo = sumDebit - sumCredt
    return kendo.toString(Saldo, 'n')
}
glidrmandiri.init = function () {
    glidrmandiri.getDataglidrmandiri(function () {
        glidrmandiri.renderGrid()
    })
}
$(document).ready(function () {
    glidrmandiri.init()
})