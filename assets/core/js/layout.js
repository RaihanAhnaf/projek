var sidebarMenu = {
	menu: ko.observableArray([])
}

sidebarMenu.getDataMenu = function () {
    ajaxPost("/menusetting/getmenutop", {}, function (data) {
    	var resMenu = data.Data.Records;
		var sortedMenu = Enumerable.From(resMenu).OrderBy("$.Parent").ThenBy("$.IndexMenu").ToArray();

		var jadi = []

		_.each(sortedMenu, function(sm){ sm.isActive = false })

		var currentMenu = _.find(sortedMenu, function(v){
			return userinfo.menunameh() == v.Title
		});
		currentMenu.isActive = true

		var setParentActive = function(){
			if(currentMenu.Parent != ""){
				currentMenu = _.find(sortedMenu, function(v){
					return currentMenu.Parent == v.Id
				});
				currentMenu.isActive = true
				setParentActive()
			} else
				return
		}
		setParentActive()

		_.each(sortedMenu, function(sm){
			if (sm.Parent == "")
				jadi.push(sm)
			else {
				var smParent = _.find(jadi, function(j){
					return j.Id == sm.Parent
				});

				if(smParent.children == undefined)
					smParent.children = []
					smParent.children.push(sm)
			}
		})
		sidebarMenu.menu(jadi)
    });
}

$(document).ready(function(){
	// kendo.bind($('#loggedinuserinfo'), userinfo)
	// kendo.bind($(".page-sidebar-menu"), sidebarMenu);
	// ko.applyBindings(userinfo)
	// ko.applyBindings(sidebarMenu)
	sidebarMenu.getDataMenu()
})