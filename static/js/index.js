function getVideoListByAjax(data) {
    //先清空元素
    $("div#query-list").empty();
    $("div#video-list").empty();

    $.get("/video-list"+data,function(data, status){
        var aa=`<a href="javascript:void(0);" onclick="getVideoListByAjax('')">全部 </a>`;
        $("div#query-list").append("请选择你的XP ",aa);
        for (var i=0;i<data.labelList.length;i++) {
            var aa=`<a href="javascript:void(0);" onclick="getVideoListByAjax('` + `?label=` + data.labelList[i] +`')">`+ data.labelList[i] +` </a>`;
            $("div#query-list").append(aa);
        }

        $("div#query-list").append("<br>");

        var aa=`<a href="javascript:void(0);" onclick="getVideoListByAjax('')">我全都要 </a>`;
        $("div#query-list").append("请选择你的LP ", aa);
        for (var i=0;i<data.performerList.length;i++) {
            var aa=`<a href="javascript:void(0);" onclick="getVideoListByAjax('` + `?performer=`  + data.performerList[i] + `')">` + data.performerList[i] +` </a>`;
            $("div#query-list").append(aa);
        }

        for (var i=0;i<data.videoList.length;i++) {
            var dd = `<div class='video-item'>` + `<a target=_blank href='/video?vid=`+ data.videoList[i].id +`'>`+ `<img class="video-cover" src="`+ data.videoList[i].coverurl +`"> <br>`+ data.videoList[i].name +`</a>` + `</div>`;
            $("div#video-list").append(dd);
        }
    });
}
