var ProActiffeThemes = {};
ProActiffeThemes.Themes = {};
ProActiffeThemes.CurrentClass = "";
ProActiffeThemes.ChangeTo = function(name) {
    for(var n in ProActiffeThemes.Themes) {
        if (!name || n == name) {
            $('html').removeClass(ProActiffeThemes.CurrentClass);
            $('html').addClass(n);
            ProActiffeThemes.CurrentClass = n;
            if(localStorage) {
                localStorage.setItem("theme", n);
            }
            if (!name) return;
        }
    }
};
$(window).on('storage', function(ev) {
    if (ev.originalEvent.key != "theme") return;
    var nTheme = ev.originalEvent.newValue;
    ProActiffeThemes.ChangeTo(nTheme);
});

 $(function() {
    /* Styling Fix for elements that are styled using js :( */
    $(".portlet-title button").first().addClass("active");
    $(".portlet-title button").each(function(i, ele) {
        var onClick = $(ele).attr("onclick");
        if (onClick) {
            if (typeof onClick != "string" || !onClick.includes("choose")) return;
            $(ele).addClass("mode-choose");
            $(ele).attr("onclick", "");
            $(ele).attr("click-callback", onClick);
            $(ele).on("click", function() {
                $(".portlet-title button").removeClass("active").attr("style", "");
                $(this).addClass("active");
                var ev = $(this).attr("click-callback");
                if (ev) $('<button onclick="' + ev + '"></button>').trigger("click");
            });
        }
    });
    $(".custom-theme").each(function(){
        var name = $(this).attr('data-theme-name');
        var color = $(this).attr('data-theme-color');
        if (name && color) {
            ProActiffeThemes.Themes[name] = color;
        }
    });

    var loadTheme = localStorage ? localStorage.getItem("theme") : "";
    ProActiffeThemes.ChangeTo(loadTheme);

    $("#themeselector").each(function() {
        for(var name in ProActiffeThemes.Themes) {
            var html = '<div alt="' + name + '" title="' + name + '" onClick="ProActiffeThemes.ChangeTo(\'' + name + '\')" style="background-color: ' + ProActiffeThemes.Themes[name] + '"></div>';
            $(this).append(html);
        }
    });
});