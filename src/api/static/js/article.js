var UpdateRead = function (articleID) {
    target_article = document.getElementById("read_"+articleID);
    if (target_article.value == "既読にする") {
        fetch('/'+articleID+'/read?read=1').then(function (response) {
            return response.text();
        }).then(function (text) {
            if (text == ''){
                target_article.value = "未読にする"
            }
        });
    } else {
        fetch('/'+articleID+'/read?read=0').then(function (response) {
            return response.text();
        }).then(function (text) {
            if (text == ''){
                target_article.value = "既読にする"
            }
        });
    }
}

var updatePriority = function (articleID, priority) {
    target_article = document.getElementById("priority_"+articleID+"_"+priority);
    fetch('/'+articleID+'/priority?priority='+priority).then(function (response) {
        return response.text();
    }).then(function (text) {
        if (text == ''){
            for(var i=1; i<=priority; i++){
                target_article = document.getElementById("priority_"+articleID+"_"+i);
                target_article.innerHTML = '<i class="fas fa-star" style="color:rgb(253,171,10);"></i>';
            }
            for(var i=priority+1; i<=5; i++){
                target_article = document.getElementById("priority_"+articleID+"_"+i);
                target_article.innerHTML = '<i class="far fa-star" style="color:rgb(253,171,10);"></i>';
            }
        }
    });
}