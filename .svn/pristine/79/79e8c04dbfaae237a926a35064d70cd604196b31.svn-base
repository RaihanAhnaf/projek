var balancesheet = {}
balancesheet.DateStart = ko.observable(moment().startOf('month').format('YYYY-MM-DD hh:mm'))
balancesheet.DateEnd = ko.observable(new Date)
balancesheet.dataMaster = ko.observableArray([])
balancesheet.dataCurrentAsset = ko.observableArray([])
balancesheet.sumDataCurrentAsset = ko.observable(0)
balancesheet.dataAktivaLain = ko.observableArray([])
balancesheet.sumdataAktivaLain = ko.observable(0)
balancesheet.dataCurrentLiabilities = ko.observableArray([])
balancesheet.sumDataCurrentLiabilities = ko.observable(0)
balancesheet.dataFixAssets = ko.observableArray([])
balancesheet.sumDataFixAssetsDebet = ko.observable(0)
balancesheet.sumDataFixAssets = ko.observable(0)
balancesheet.totalAssets = ko.observable(0)
balancesheet.dataCapital = ko.observable(0)
balancesheet.dataRetained = ko.observable(0)
balancesheet.dataDeviden = ko.observable(0)
balancesheet.dataCurrentEarning = ko.observable(0)
balancesheet.sumCapitalandEarning = ko.observable(0)
balancesheet.totalpassiva = ko.observable(0)
balancesheet.paramFilter = ko.observable(false)
balancesheet.accumulatedDepreciation = ko.observable(0)
balancesheet.resultspasiva = ko.observableArray([])
balancesheet.resultsactiva = ko.observableArray([])
balancesheet.capitalAndEarningRow = ko.observableArray([])
balancesheet.fixAssetRow = ko.observableArray([])
balancesheet.TitelFilter = ko.observable(" Hide Filter")
balancesheet.DatePageBar = ko.observable()
balancesheet.ShowMonthFilter = ko.observable()

// Fix Assets
balancesheet.activaTetap = ko.observableArray([])
balancesheet.building = ko.observableArray([])
balancesheet.accDepBuilding = ko.observableArray([])
balancesheet.vehicle = ko.observableArray([])
balancesheet.accDepVehicle = ko.observableArray([])
balancesheet.officeEquipment = ko.observableArray([])
balancesheet.accDepOfficeEquipment = ko.observableArray([])
balancesheet.aktivaLainnya = ko.observableArray([])



balancesheet.getDateNow = function() {
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    balancesheet.DatePageBar(page)
}

balancesheet.getDataCurrentAsset = function(callback) {
    model.Processing(true)
    var url = "/financial/getdatabalancesheetcoa"
    var param = {
        DateStart: $('#dateStartBalance').data('kendoDatePicker').value(),
        DateEnd: $('#dateEndBalance').data('kendoDatePicker').value(),
        Category: "BALANCE SHEET",
        DebCred: "DEBET",
        Filter: balancesheet.paramFilter(),
    }
    ajaxPost(url, param, function(res) {
        model.Processing(false)
        //Current Assets
        res.Data = _.sortBy(res.Data, [function(o) { return o.ACC_Code; }]);
        var c = _.filter(res.Data, function(o) {
            if (o.ACC_Code==1401){
                o.Account_Name = "PERSEDIAAN"
            }
            return o.Main_Acc_Code > 0 && o.ACC_Code < 2000 ;
        });
        // _.sortBy(c, [function(o) { return o.ACC_Code; }]);
        var sumc = _.sumBy(c, 'Saldo')
        var aktvaLainnya = _.filter(res.Data, function(o) {
            return o.ACC_Code > 2500 && o.ACC_Code < 3000;
        });
        var sumaktvaLainnya = _.sumBy(aktvaLainnya, 'Saldo')
        //Fix Asset Debit
        // building
        // var building = _.filter(res.Data, function(o) {
        //     return o.ACC_Code == 2200;
        // });
    

        // // vehicle
        // var vehicle = _.filter(res.Data, function(o) {
        //     return o.ACC_Code == 2300;
        // });
        var activatetap = _.filter(res.Data, function(o) {
            return o.ACC_Code > 2000 && o.ACC_Code< 2500;
        });
        //aktiva lainnya 
        var activaLainnya = _.filter(res.Data, function(o) {
            return o.ACC_Code > 2500&& o.ACC_Code< 3000;;
        });

        // office equipment -- INVENTARIS KANTOR
        var office = _.filter(res.Data, function(o) {
            return o.ACC_Code == 2400;
        });
        // Accumulated Depreciation Building
        var adbuilding = _.filter(res.Data, function(o) {
            return o.ACC_Code == 2210; //2600 --> 2210
        });
        // Accumulated Depreciation Vehicle
        var advehicle = _.filter(res.Data, function(o) {
            return o.ACC_Code == 2310; //2700 --> 2310
        });

        //Accumulated Depreciation Office Equipment
        var adoffice = _.filter(res.Data, function(o) {
            return o.ACC_Code == 2410; //2800 --> 2410
        });
        console.log("--",adbuilding)
        // if (c.length > 0) {
            balancesheet.dataCurrentAsset(c)
            balancesheet.sumDataCurrentAsset(sumc)
            balancesheet.dataAktivaLain(aktvaLainnya)
            balancesheet.sumdataAktivaLain(sumaktvaLainnya)
        // }
        if (adbuilding.length > 0) {
            balancesheet.accDepBuilding(adbuilding[0].Saldo)
            $("#balancesheet-accDepBuilding").show();
        } else {
            balancesheet.accDepBuilding([])
            $("#balancesheet-accDepBuilding").hide();
            // balancesheet.accDepBuilding().Saldo = 0
        } 

        if (advehicle.length > 0) {
            balancesheet.accDepVehicle(advehicle[0].Saldo)
            $("#balancesheet-advehicle").show();
        } else {
            balancesheet.accDepVehicle([])
            $("#balancesheet-advehicle").hide();
            // balancesheet.accDepVehicle().Saldo = 0
        }

        if (adoffice.length > 0) {
            balancesheet.accDepOfficeEquipment(adoffice[0])
            balancesheet.accumulatedDepreciation(adoffice[0].Saldo)
            $("#balancesheet-accDepOfficeEquipment").show();
            $("#balancesheet-accumulatedDepreciation").show();
        } else {
            balancesheet.accDepOfficeEquipment([])
            balancesheet.accumulatedDepreciation([])
            $("#balancesheet-accDepOfficeEquipment").hide();
            $("#balancesheet-accumulatedDepreciation").hide();
            // balancesheet.accDepOfficeEquipment().Saldo = 0
        }

        if (activatetap.length > 0) {
            balancesheet.activaTetap(activatetap);
            $("#balancesheet-building").show();
        } else {
            balancesheet.activaTetap([]);
            $("#balancesheet-building").hide();
            // balancesheet.building().Saldo = 0
        }
        // if (vehicle.length > 0) {
        //     balancesheet.vehicle(vehicle[0].Saldo)
        //     $("#balancesheet-vehicle").show();
        // } else {
        //     balancesheet.vehicle([])
        //     $("#balancesheet-vehicle").hide();
        //     // balancesheet.vehicle().Saldo = 0
        // }
        if (office.length > 0) {
            balancesheet.officeEquipment(office[0])
            $("#balancesheet-officeEquipment").show();
        } else {
            balancesheet.officeEquipment([])
             $("#balancesheet-officeEquipment").hide();
            // balancesheet.officeEquipment().Saldo = 0
        }
        if (activaLainnya.length > 0) {
            balancesheet.aktivaLainnya(activaLainnya[0].Saldo)
            $("#balancesheet-aktivaLainnya").show();
        } else {
            balancesheet.aktivaLainnya([])
            $("#balancesheet-aktivaLainnya").hide();
            // balancesheet.vehicle().Saldo = 0
        }

        callback()
    }, function() {
        swal("Error!", "Unknown error, please try again", "error")
    })
}

balancesheet.getDataCurrentLiabilities = function() {
    var url = "/financial/getdatabalancesheetcoa"
    var param = {
        DateStart: $('#dateStartBalance').data('kendoDatePicker').value(),
        DateEnd: $('#dateEndBalance').data('kendoDatePicker').value(),
        Category: "BALANCE SHEET",
        DebCred: "CREDIT",
        Filter: balancesheet.paramFilter(),
    }
    ajaxPost(url, param, function(res) {
        res.Data = _.sortBy(res.Data, [function(o) { return o.ACC_Code; }]);
        //Current Liabilities
        var c = _.filter(res.Data, function(o) {
            return o.Main_Acc_Code > 0 && o.ACC_Code < 4000 && o.ACC_Code > 3000;
        });
        var sumc = _.sumBy(c, 'Saldo')

        //Capital
        var d = _.filter(res.Data, function(o) {
            return o.ACC_Code == 4100;
        });
        var sumd = _.sumBy(d, 'Saldo')


        // Fix Asset Kredit


        //Retained Earning
        var f = _.filter(res.Data, function(o) {
            return o.ACC_Code == 4200;
        });
         //Retained Earning
        var h = _.filter(res.Data, function(o) {
            return o.ACC_Code == 4300;
        });

        //Current Earning
        var g = _.filter(res.Data, function(o) {
            return o.ACC_Code == 4400;
        });

        // if (c.length > 0) {
             balancesheet.dataCurrentLiabilities(c) 
        // }
        if (d.length > 0) {
              balancesheet.dataCapital(d[0].Saldo)
              $("#balancesheet-dataCapital").show();
        }else {
            balancesheet.dataCapital(0)
            $("#balancesheet-dataCapital").hide();
        }
        if (f.length > 0) {
            balancesheet.dataRetained(f[0].Saldo)
            $("#balancesheet-dataRetained").show();
        }else{
            balancesheet.dataRetained(0)
            $("#balancesheet-dataRetained").hide();
        }
        if (h.length > 0) {
            balancesheet.dataDeviden(h[0].Saldo)
            $("#balancesheet-dataDeviden").show();
        }else{
            balancesheet.dataDeviden(0)
            $("#balancesheet-dataDeviden").hide();
        }
        if (g.length > 0) {
             balancesheet.dataCurrentEarning(g[0].Saldo)
             $("#balancesheet-dataCurrentEarning").show();
        }else{
            balancesheet.dataCurrentEarning(0)
            $("#balancesheet-dataCurrentEarning").hide();
        }

        // var total = (balancesheet.building() + balancesheet.vehicle() + balancesheet.officeEquipment().Saldo) - (balancesheet.accDepBuilding() + balancesheet.accDepVehicle() + balancesheet.accDepOfficeEquipment().Saldo) * -1
        var total = _.sumBy(balancesheet.activaTetap(), function(e){return e.Saldo})
        if (total > 0){
            balancesheet.sumDataFixAssets(total)
        }

        var totalAsset = balancesheet.sumDataCurrentAsset() + balancesheet.sumDataFixAssets()+ balancesheet.sumdataAktivaLain()
        balancesheet.totalAssets(totalAsset)
        balancesheet.sumDataCurrentLiabilities(sumc)

        var totalcapitalandearning = balancesheet.dataCapital() + (balancesheet.dataRetained()+balancesheet.dataDeviden() + balancesheet.dataCurrentEarning())
        balancesheet.sumCapitalandEarning(totalcapitalandearning)

        var totalpassiva = balancesheet.sumDataCurrentLiabilities() + balancesheet.sumCapitalandEarning()
        balancesheet.totalpassiva(totalpassiva)
        balancesheet.rowEmpty()
    })
}

balancesheet.refresh = function() {
    balancesheet.resultspasiva([])
    balancesheet.resultsactiva([])
    balancesheet.capitalAndEarningRow([])
    balancesheet.fixAssetRow([])
    balancesheet.paramFilter(true)
    balancesheet.getDataCurrentAsset(function() {
        balancesheet.getDataCurrentLiabilities()
    })
    var from = kendo.toString($('#dateStartBalance').data('kendoDatePicker').value(), 'MMMM yyyy')
    var to = kendo.toString($('#dateEndBalance').data('kendoDatePicker').value(), 'MMMM yyyy')
    if (from == to) {
        balancesheet.ShowMonthFilter(from)
    }else{
        balancesheet.ShowMonthFilter(from+ " - "+to)
    }
}

balancesheet.setDate = function() {
    var datepicker = $("#dateStartBalance").data("kendoDatePicker");
    var now = new Date();
    var years = moment(now).format("YYYY")
    var Month = moment(now).format("MM")
    var newDate = years + "-" + Month + "-" + "01"
    datepicker.value(new Date(newDate))
    balancesheet.ShowMonthFilter(moment(now).format("MMMM YYYY"))
}

balancesheet.ExportToPdfBalanceSheet = function() {
    model.Processing(true)
    var url = "/financial/exporttopdfbalancesheet"
    var param = {
        DateStart: $('#dateStartBalance').data('kendoDatePicker').value(),
        DateEnd: $('#dateEndBalance').data('kendoDatePicker').value(),
        Filter: balancesheet.paramFilter(),
    }
    ajaxPost(url, param, function(res) {
        model.Processing(false)
        window.open('/res/docs/balancesheet/' + res, '_blank');
    })
}

balancesheet.rowEmpty = function() {
    balancesheet.toggleFilter = function() {
        var panelFilter = $('.panel-filter');
        var panelContent = $('.panel-content');

        if (panelFilter.is(':visible')) {
            panelFilter.hide();
            panelContent.attr('class', 'col-md-12 col-sm-12 ez panel-content');
            $('.breakdown-filter').removeAttr('style');
        } else {
            panelFilter.show();
            panelContent.attr('class', 'col-md-9 col-sm-9 ez panel-content');
            $('.breakdown-filter').css('width', '60%');
        }

        $('.k-grid').each(function(i, d) {
            try {
                $(d).data('kendoGrid').refresh();
            } catch (err) {}
        });

        $('.k-pivot').each(function(i, d) {
            $(d).data('kendoPivotGrid').refresh();
        });

        $('.k-chart').each(function(i, d) {
            $(d).data('kendoChart').redraw();
        });
        balancesheet.panel_relocated();
        var FilterTitle = balancesheet.TitelFilter();
        if (FilterTitle == " Hide Filter") {
            balancesheet.TitelFilter(" Show Filter");
        } else if (FilterTitle == " Show Filter") {
            balancesheet.TitelFilter(" Hide Filter");
        }
    }
    balancesheet.panel_relocated = function() {
        if ($('.panel-yo').size() == 0) {
            return;
        }
        var window_top = $(window).scrollTop();
        var div_top = $('.panel-yo').offset().top;
        if (window_top > div_top) {
            $('.panel-fix').css('width', $('.panel-yo').width());
            $('.panel-fix').addClass('contentfilter');
            $('.panel-yo').height($('.panel-fix').outerHeight());
        } else {
            $('.panel-fix').removeClass('contentfilter');
            $('.panel-yo').height(0);
        }
    }
    var x = 0
    var a = balancesheet.dataCurrentAsset().length
    var b = balancesheet.dataCurrentLiabilities().length

    if (a > b) {
        x = a - b
        for (i = 0; i < x; i++) {
            balancesheet.resultspasiva.push({
                "data": i
            })
        }
    } else {
        x = b - a
        for (i = 0; i < x; i++) {
            balancesheet.resultsactiva.push({
                "data": i
            })
        }
    }
    var a2 = $('#fix-asset-tbody tr:visible').length
    var b2 = $('#capital-and-earning tr:visible').length
    if (a2 > b2) {
        xx = a2 - b2
        for (i = 0; i < xx; i++) {
            balancesheet.capitalAndEarningRow.push({
                "data": i
            })
        }
    } else {
        xx = b2 - a2
        for (i = 0; i < xx; i++) {
            balancesheet.fixAssetRow.push({
                "data": i
            })
        }
    }
}

balancesheet.onChangeDateStart = function(val){
    if (val.getTime()>balancesheet.DateEnd().getTime()){
        balancesheet.DateEnd(val)
    }
}

balancesheet.init = function() {
    balancesheet.getDateNow()
    balancesheet.ShowMonthFilter(moment().format("MMMM YYYY"))
    // balancesheet.setDate()
    balancesheet.getDataCurrentAsset(function() {
        balancesheet.getDataCurrentLiabilities()
    })
}


$(function() {
    balancesheet.init()
})