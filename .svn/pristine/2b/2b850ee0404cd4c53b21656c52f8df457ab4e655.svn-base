model.Processing(false)

var transactionorder = {}
transactionorder.DatePageBar = ko.observable()

transactionorder.getDateNow = function() {
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    transactionorder.DatePageBar(page)
}

transactionorder.newRecord = function() {
    var page = {
        ID: "",
        DateStr: "",
        DatePosting: "",
        DocumentNumber: "",
        StoreHouseFrom: "",
        StoreHouseTo: "",
        Description: "",
        ListDetailOrder: [],
    }
    page.ListDetailOrder.push(transactionorder.listDetailOrder({}))
    return page
}

transactionorder.listDetailOrder = function(data) {
    var dataTmp = {}
    dataTmp.Id = data.Id == undefined ? '' : data.Id
    dataTmp.Item = data.Item == undefined ? '' : data.Item    
    dataTmp.StockUnit = data.StockUnit == undefined ? '' : data.StockUnit
    dataTmp.Qty = data.Qty == undefined ? '' : data.Qty
    dataTmp.CodeItem = data.CodeItem == undefined ? '' : data.CodeItem
    var x = ko.mapping.fromJS(dataTmp)

    return x
}

transactionorder.record = ko.mapping.fromJS(transactionorder.newRecord())
transactionorder.createdForm = function() {

}

transactionorder.warehouseFrom = ko.observableArray([]);
transactionorder.warehouseTo = ko.observableArray([]);
transactionorder.getMasterLocation = function() {
    $.ajax({
        url: "/transaction/getdatamasterlocation",
        success: function(json) {
            if (!json.IsError) {
                transactionorder.warehouse = [];
                $(json.Data).each(function(ix, ele) {
                    transactionorder.warehouseFrom.push({
                        value: ele._id,
                        text: ele.LocationName
                    });
                    transactionorder.warehouseTo.push({
                        value: ele._id,
                        text: ele.LocationName
                    });
                });
            }
        }
    });
}

transactionorder.init = function() {
    transactionorder.getDateNow();
    transactionorder.getMasterLocation();
}

$(function() {
transactionorder.init()
});