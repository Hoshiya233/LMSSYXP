function getVideoListByAjax(filter_str) {
    //先清空元素
    $("div#filter-list").empty();
    $("div#video-list").empty();

    $.get("/video-list"+filter_str,function(data, status){
        var aa=`<button class="btn-default" href="javascript:void(0);" onclick="getVideoListByAjax('')">全部</button>`;
        $("div#filter-list").append("请选择你的XP ",aa);
        for (var i=0;i<data.labelList.length;i++) {
            var aa=`<button class="btn-default" href="javascript:void(0);" onclick="getVideoListByAjax('` + `?label=` + data.labelList[i] +`')">`+ data.labelList[i] +`</button>`;
            $("div#filter-list").append(aa);
        }

        $("div#filter-list").append("<br>");

        var aa=`<button class="btn-default" href="javascript:void(0);" onclick="getVideoListByAjax('')">我全都要</button>`;
        $("div#filter-list").append("请选择你的LP ", aa);
        for (var i=0;i<data.performerList.length;i++) {
            var aa=`<button class="btn-default" href="javascript:void(0);" onclick="getVideoListByAjax('` + `?performer=`  + data.performerList[i] + `')">` + data.performerList[i] +`</button>`;
            $("div#filter-list").append(aa);
        }

        for (var i=0;i<data.videoList.length;i++) {
            var dd = `<div class="video-item">` + `<a class="video-link" target="_blank" href="/video?vid=`+ data.videoList[i].id +`">`+ `<img class="video-cover" src="`+ data.videoList[i].coverurl +`"> <br>`+ data.videoList[i].name +`</a>` + `</div>`;
            $("div#video-list").append(dd);
        }
    });
}
