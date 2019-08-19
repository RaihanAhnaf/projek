// var menusett = kendo.observable({
// 	Id: "",
// 	PageId: "",
//     Parent: "",
//     Title: "",
//     Url: "",
//     Icon: "",
//     IndexMenu: 0,
//     saveData : false,
//     updateData : false,
//     DataMenu : [],
//     listMenu : [],
//     ListMenuTree : [],
//     treelistView : [],
//     btnNewData : " Add New",
//     select: false,
// });
var menusett ={
	Id: ko.observable(""),
	PageId: ko.observable(""),
    Parent: ko.observable(""),
    Title: ko.observable(""),
    Url: ko.observable(""),
    Icon: ko.observable(""),
    IndexMenu: ko.observable(0),
    saveData : ko.observable(false),
    updateData : ko.observable(false),
    DataMenu : ko.observableArray([]),
    listMenu : ko.observableArray([]),
    ListMenuTree : ko.observableArray([]),
    treelistView : ko.observableArray([]),
    btnNewData : ko.observable(" Add New"),
    select: ko.observable(false),
    DatePageBar: ko.observable()
};
// var ddparentVM = kendo.observable({
// 	Parent: null, 	//Parent
// 	listMenu: [], 	//listMenu
// });
var ddparentVM = {
	Parent: ko.observable(null), 	//Parent
	listMenu: ko.observableArray([]), 	//listMenu
};

menusett.getDateNow = function(){
    var page = moment(new Date()).format("DD MMMM YYYY HH:mm")
    menusett.DatePageBar(page)
}

menusett.resetAppMenu = function(){
	menusett.Id("");
    menusett.PageId("");
    ddparentVM.Parent("");
    menusett.Title("");
    menusett.Url("");
    menusett.Icon("");
    menusett.IndexMenu(0);
    menusett.saveData(true);
    menusett.updateData(false);
    // menusett.set("Id", "");
    // menusett.set("PageId", "");
    // ddparentVM.set("Parent", "");
    // menusett.set("Title", "");
    // menusett.set("Url", "");
    // menusett.set("Icon", "");
    // menusett.set("IndexMenu", 0);
    // menusett.set("saveData", true);
    // menusett.set("updateData", false);

    $('#Enable').bootstrapSwitch('state', true)
    $("#Url").siblings("span.k-tooltip-validation").hide(); 
    $("#title").siblings("span.k-tooltip-validation").hide(); 

    menusett.loadmenuMaster();
}

menusett.saveAppmenu = function(){
	// var Title = menusett.get("Title");
	var Title = menusett.Title();
	var IndexMenu = menusett.IndexMenu();
	// var IndexMenu = menusett.get("IndexMenu");
	var Enable = $('#Enable').bootstrapSwitch('state');
	var d = new Date();
    var day = d.getDate();
    var month = d.getMonth() + 1;
    var year = d.getFullYear();
    var Hours= d.getHours();
    var Minutes= d.getMinutes();
    var Seconds= d.getSeconds();
    var GenMenuId = year + "" + month + "" + day + "" + Hours + "" + Minutes + "" + Seconds

	if (IndexMenu == ""){
		IndexMenu = 0;
	}

	var validator = $("#AppMenu").data("kendoValidator");

    if (validator == undefined){
       validator= $("#AppMenu").kendoValidator().data("kendoValidator");
    }

     if (validator.validate()) {
		var param = {
			Id : GenMenuId,
			PageId : Title.toUpperCase().replace(/\s+/g, ''),
			Parent : ddparentVM.Parent(),
			// Parent : ddparentVM.get("Parent"),
			Title : Title,
			Url : menusett.Url(),
			Icon : menusett.Icon(),
			// Url : menusett.get("Url"),
			// Icon : menusett.get("Icon"),
			IndexMenu : parseInt(IndexMenu),
			Enable : Enable
		};

     	ajaxPost("/menusetting/savemenutop", param, function (data) {
			if (data.IsError == false){
			    menusett.resetAppMenu();

			    // menusett.get("btnNewData", " Add New");
			    // menusett.get("saveData", false);
			    menusett.btnNewData(" Add New");
			    menusett.saveData(false);

				$("#parent").data("kendoDropDownList").readonly(true);
				$("#IndexMenu").data("kendoNumericTextBox").readonly(true);
				$("#Url").attr("readonly", true);
				$("#title").attr("readonly", true);
				$("#Icon").attr("readonly", true);
				$('#Enable').bootstrapSwitch('readonly', true);
				swal({
					title:"Success!",
					text: data.Message,
					type: "success",
					confirmButtonColor: "#3da09a"});
			}else{
				return swal({
					title:"Error!",
					text: data.Message,
					type: "error",
					confirmButtonColor:"#3da09a"});
			}
		});
     }
	
}

menusett.updateAppmenu = function(){
	// var Title = menusett.get("Title");
	// var IndexMenu = menusett.get("IndexMenu");
	var Title = menusett.Title();
	var IndexMenu = menusett.IndexMenu();
	var Enable = $('#Enable').bootstrapSwitch('state');

	if (IndexMenu == ""){
		IndexMenu = 0;
	}

	var param = {
		// Id : menusett.get("Id"),
		// PageId : menusett.get("PageId"),
		// Parent : ddparentVM.get("Parent"),
		Id : menusett.Id(),
		PageId : menusett.PageId(),
		Parent : ddparentVM.Parent(),
		Title : Title,
		// Url : menusett.get("Url"),
		// Icon : menusett.get("Icon"),
		Url : menusett.Url(),
		Icon : menusett.Icon(),
		IndexMenu : parseInt(IndexMenu),
		Enable : Enable
	};

	ajaxPost("/menusetting/updatemenutop", param, function (data) {
		if (data.IsError == false){
		    menusett.resetAppMenu();

		    // menusett.get("select", false);
		    // menusett.get("saveData", false);
		    menusett.select(false);
		    menusett.saveData(false);

		    $("#parent").data("kendoDropDownList").readonly(true);
			$("#IndexMenu").data("kendoNumericTextBox").readonly(true);
			$("#Url").attr("readonly", true);
			$("#title").attr("readonly", true);
			$("#Icon").attr("readonly", true);
			$('#Enable').bootstrapSwitch('readonly', true);

			swal({
				title: "Success!",
				text: data.Message,
				type: 'success',
				showCancelButton: false,
				confirmButtonColor: '#3da09a',
				cancelButtonColor: '#d33',
			}, function(dismiss) {
				console.log(dismiss)
				if (dismiss === 'cancel') {
					console.log("dismiss");
				} else {
					location.reload();
				}
			});
		}else{
			return swal({
				title:"Error!",
				text: data.Message,
				type: "error",
				confirmButtonColor:"#3da09a"});
		}
	});
}

menusett.checkSelect = function(){
	var tv = $("#menu-list").data("kendoTreeView");
	selected = tv.select();
	item = tv.dataItem(selected);

    if (item === undefined) {
       	return swal({
       		title:"Confirmation!",
       		text: "Please select menu.",
       		type: "error",
       		confirmButtonColor:"#3da09a"
       	});
    } else {
    	// menusett.set("Id", item.Id);
    	// menusett.set("saveData", true);
    	// menusett.set("updateData", true);
    	menusett.Id(item.Id);
    	menusett.saveData(true);
    	menusett.updateData(true);
    }
}

menusett.editdataMenulist = function(){
	menusett.checkSelect();

	$("#parent").data("kendoDropDownList").readonly(false);
	$("#IndexMenu").data("kendoNumericTextBox").readonly(false);
	$("#Url").attr("readonly", false);
	$("#title").attr("readonly", false);
	$("#Icon").attr("readonly", false);
	$('#Enable').bootstrapSwitch('readonly', false);

	// menusett.set("saveData", false);
	menusett.saveData(false);

}

menusett.newdataMenulist = function(){
	// menusett.set("select", false);
	menusett.select(false);

	$("#parent").data("kendoDropDownList").readonly(false);
	$("#IndexMenu").data("kendoNumericTextBox").readonly(false);
	$("#Url").attr("readonly", false);
	$("#title").attr("readonly", false);
	$("#Icon").attr("readonly", false);
	$('#Enable').bootstrapSwitch('readonly', false);	

	// menusett.set("saveData", true);
 //    menusett.set("updateData", false);
 //    menusett.set("Id", "");
 //    menusett.set("PageId", "");
 //    ddparentVM.set("Parent", "");
 //    menusett.set("Title", "");
 //    menusett.set("Url", "");
 //    menusett.set("Icon", "");
 //    menusett.set("IndexMenu", 0);
 	menusett.saveData(true);
    menusett.updateData(false);
    menusett.Id("");
    menusett.PageId("");
    ddparentVM.Parent("");
    menusett.Title("");
    menusett.Url("");
    menusett.Icon("");
    menusett.IndexMenu(0);

    menusett.loadmenuMaster();

    $('#Enable').bootstrapSwitch('state', true)
    $("#Url").siblings("span.k-tooltip-validation").hide(); 
    $("#title").siblings("span.k-tooltip-validation").hide(); 

    // if (menusett.get("btnNewData") == " Add New"){
    // 	menusett.set("btnNewData", " Cancel");
    // 	menusett.set("saveData", true);
    // }else{
    // 	menusett.set("btnNewData", " Add New");
    // 	menusett.set("saveData", false);
    // 	menusett.set("updateData", false);
    if (menusett.btnNewData()== " Add New"){
    	menusett.btnNewData(" Cancel");
    	menusett.saveData(true);
    }else{
    	menusett.btnNewData(" Add New");
    	menusett.saveData(false);
    	menusett.updateData(false);

    	$("#parent").data("kendoDropDownList").readonly(true);
		$("#IndexMenu").data("kendoNumericTextBox").readonly(true);
		$("#Url").attr("readonly", true);
		$("#title").attr("readonly", true);
		$("#Icon").attr("readonly", true);
		$('#Enable').bootstrapSwitch('readonly', true);
    }
}


menusett.deleteMunulist = function(){
	// if (menusett.get("Id") == ""){
	if (menusett.Id() == ""){
		return swal({
			title:"Confirmation!",
			text: "Please select menu.",
			type: "error",
			confirmButtonColor:"#3da09a"
		});
	}

	swal({
            title: "Are you sure?",
            text: "Are you sure remove this menu!",
            type: "warning",
            showCancelButton: true,
            confirmButtonColor: '#3da09a',
            confirmButtonText: 'Yes, I am sure!',
            cancelButtonText: "No, cancel it!",
            closeOnConfirm: false,
            closeOnCancel: false
        },
        function(isConfirm) {
            if (isConfirm) {
			    ajaxPost("/menusetting/deletemenutop", { Id : menusett.Id() }, function(data){
			    // ajaxPost("/menusetting/deletemenutop", { Id : menusett.get("Id") }, function(data){
					if (data.IsError == false){
						menusett.resetAppMenu();

						// swal("Success!", "Menu Success " + menusett.get("Id") + " Delete", "success");
						// menusett.set("saveData", false);
						swal({
							title:"Success!",
							text: "Menu Success " + menusett.Id() + " Delete",
							type: "success",
							confirmButtonColor:"#3da09a"});
						menusett.saveData(false);
					}else{
						swal("Error!",data.Message,"error");
						// menusett.set("saveData", false);
						menusett.saveData(false);
					}
				});
            } else {
            	menusett.saveData(false);
			    menusett.updateData(false);
			    menusett.Id("");
			    menusett.PageId("");
			    ddparentVM.Parent("");
			    menusett.Title("");
			    menusett.Url("");
			    menusett.Icon("");
			    menusett.IndexMenu(0);
       			// menusett.set("saveData", false);
			    // menusett.set("updateData", false);
			    // menusett.set("Id", "");
			    // menusett.set("PageId", "");
			    // ddparentVM.set("Parent", "");
			    // menusett.set("Title", "");
			    // menusett.set("Url", "");
			    // menusett.set("Icon", "");
			    // menusett.set("IndexMenu", 0);

			    menusett.loadmenuMaster();

			    $('#Enable').bootstrapSwitch('state', true)
			    $("#Url").siblings("span.k-tooltip-validation").hide(); 
			    $("#title").siblings("span.k-tooltip-validation").hide(); 

                swal({
                	title:"Cancelled",
                	text: "Cancelled Delete Menu",
                	type: "error",
                	confirmButtonColor:"#3da09a"});
            }
        });
}

menusett.convert = function (array){
    var map = {};
    for(var i = 0; i < array.length; i++){
        var obj = array[i];
        obj.Submenus= [];

        map[obj.Id] = obj;

        var parent = obj.Parent || '-';
        if(!map[parent]){
            map[parent] = {
                Submenus: []
            };
        }
        map[parent].Submenus.push(obj);
    }
    return map['-'].Submenus;
}


menusett.subMenuMaster = function(SubData, spacer){
	spacer += "--";
	for (var i in SubData){
			if (SubData[i].Submenus.length != 0 ){
				ddparentVM.listMenu.push({
						"title" : spacer + " " + SubData[i].Title,
						"Id" : SubData[i].Id
					});
				menusett.subMenuMaster(SubData[i].Submenus, spacer);
			}else{
				ddparentVM.listMenu.push({
						"title" : spacer + " " + SubData[i].Title,
						"Id" : SubData[i].Id
					});
			}
		}
}

menusett.subtreelist = function(SubData){
	for (var i in SubData){
			if (SubData[i].Submenus.length != 0 ){
				menusett.treelistView.push({
						"Id" : SubData[i].Id,
						"title" : SubData[i].Title,
						"url" : "#",
						"icon" : SubData[i].Icon,
						"pageid" : SubData[i].PageId,
						"Parent" : SubData[i].Parent,
						"IndexMenu" : SubData[i].IndexMenu,
						"enable" : SubData[i].Enable,
					});
				menusett.subtreelist(SubData[i].Submenus);
			}else{
				menusett.treelistView.push({
						"Id" : SubData[i].Id,
						"title" : SubData[i].Title,
						"url" : "#",
						"icon" : SubData[i].Icon,
						"pageid" : SubData[i].PageId,
						"Parent" : SubData[i].Parent,
						"IndexMenu" : SubData[i].IndexMenu,
						"enable" : SubData[i].Enable,
					});
			}
		}
}

menusett.loadmenuMaster = function(){
	// ddparentVM.set("listMenu", [{ title: "[TOP LEVEL]", Id: "" }]);
	// menusett.set("treelistView", []);
	ddparentVM.listMenu([{ title: "[TOP LEVEL]", Id: "" }]);
	menusett.treelistView([]);

	ajaxPost("/menusetting/getselectmenu", {}, function (data) {
		var dataMenu = data.Data.Records;
		var sortdataMenu = Enumerable.From(dataMenu).OrderBy("$.Parent").ThenBy("$.IndexMenu").ToArray();
		var dataTree =  menusett.convert(sortdataMenu);
		var spacer = "--";
		var listSubmenu = [];

		for (var i in dataTree){
			if (dataTree[i].Submenus.length  != 0){
				ddparentVM.listMenu.push({
					"title" : spacer + " " + dataTree[i].Title,
					"Id" : dataTree[i].Id
				});

				menusett.subMenuMaster(dataTree[i].Submenus, spacer);

				//=================== 
				menusett.treelistView.push({
					"Id" : dataTree[i].Id,
					"title" : dataTree[i].Title,
					"url" : "#",
					"icon" : dataTree[i].Icon,
					"pageid" : dataTree[i].PageId,
					"Parent" : dataTree[i].Parent,
					"IndexMenu" : dataTree[i].IndexMenu,
					"enable" : dataTree[i].Enable,
				});

				menusett.subtreelist(dataTree[i].Submenus);
			} else {
				ddparentVM.listMenu.push({
					"title" : spacer + " " + dataTree[i].Title,
					"Id" : dataTree[i].Id
				});

				menusett.treelistView.push({
					"Id" : dataTree[i].Id,
					"title" : dataTree[i].Title,
					"url" : "#",
					"icon" : dataTree[i].Icon,
					"pageid" : dataTree[i].PageId,
					"Parent" : dataTree[i].Parent,
					"IndexMenu" : dataTree[i].IndexMenu,
					"enable" : dataTree[i].Enable,
				});

			}
		}

		// var sortdataTree =  menusett.convert(menusett.get("treelistView"));
		var sortdataTree =  menusett.convert(menusett.treelistView());
		// menusett.set("ListMenuTree", sortdataTree);
		menusett.ListMenuTree(sortdataTree);


		var inline = new kendo.data.HierarchicalDataSource({
            // data: menusett.get("ListMenuTree"),
            data: menusett.ListMenuTree(),
            schema: {
                model: {
                    children: "Submenus"
                }
            }
        });

	 	var treeview = $("#menu-list").kendoTreeView({
            animation: false,
            template: kendo.template($("#menulist-template").html()),
            dataTextField: "title",
            // dataValueField: "_id",
            dataSource:inline,
            select: menusett.selectDirFolder,
            loadOnDemand: false
        }).data("kendoTreeView");
        
        treeview.expand(".k-item");
	});
}

menusett.oncancel = function(){
	// if(menusett.get("btnNewData") == ' Cancel'){
	// 	menusett.set("saveData", false);
	if(menusett.btnNewData() == ' Cancel'){
		menusett.saveData(false);
		$("#parent").data("kendoDropDownList").readonly(true);
		$("#IndexMenu").data("kendoNumericTextBox").readonly(true);
		$("#Url").attr("readonly", true);
		$("#title").attr("readonly", true);
		$("#Icon").attr("readonly", true);
		$('#Enable').bootstrapSwitch('readonly', true);
	}	
}    

menusett.selectDirFolder = function(e){
	menusett.oncancel();

	var data = $('#menu-list').data('kendoTreeView').dataItem(e.node);

	ajaxPost("/menusetting/getselectmenu", { Id : data.Id }, function(data){
		if (data.IsError == false){
			// menusett.set("select", true);
			// menusett.set("saveData", false);
			menusett.select(true);
			menusett.saveData(false);

			var dataMenu =  data.Data.Records[0];
			// menusett.set("Id", dataMenu.Id);
			// menusett.set("PageId", dataMenu.PageId);
			// ddparentVM.set("Parent", dataMenu.Parent);
		 //    menusett.set("Title", dataMenu.Title);
		 //    menusett.set("Url", dataMenu.Url);
		 //    menusett.set("Icon", dataMenu.Icon);
		 //    menusett.set("IndexMenu", dataMenu.IndexMenu);
			menusett.Id(dataMenu.Id);
			menusett.PageId(dataMenu.PageId);
			ddparentVM.Parent(dataMenu.Parent);
		    menusett.Title(dataMenu.Title);
		    menusett.Url(dataMenu.Url);
		    menusett.Icon(dataMenu.Icon);
		    menusett.IndexMenu(dataMenu.IndexMenu);

		    $('#Enable').bootstrapSwitch('state', dataMenu.Enable)
		}else{
			swal("Error!",data.Message,"error");
		}
	});
}

$(document).ready(function () {
	menusett.getDateNow()
	setTimeout(function() {
		// style for input Menu Entry
		$("span.k-dropdown>span").css("height","32px");
		$("span.k-select>span>span.k-i-arrow-60-up").css("padding-top","6px");
	},50);
	
	// kendo.bind($("#menusett"), menusett)
	// kendo.bind($("#parent"), ddparentVM)

	// $("#parent").kendoDropDownList({
	// 	dataValueField: 'Id', 
	// 	dataTextField: 'title',
	// 	valuePrimitive: true,
	// 	optionLabel: "Select parent..."
	// })

	$("#IndexMenu").kendoNumericTextBox({
		value: menusett.IndexMenu, 
		min: 0
	})

	menusett.loadmenuMaster();
	$("#IndexMenu").data("kendoNumericTextBox").readonly();

});