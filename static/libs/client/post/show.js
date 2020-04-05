$(document).ready(function () {
    var body = $("#body");
    // markdown parse
    var md = window.markdownit({
        html: true
    });
    var result = md.render(body.text());
    body.html(result);
    $(window).scroll(function () {
        if ($(this).scrollTop() > 100) {
            $('#back-to-top').fadeIn();
        } else {
            $('#back-to-top').fadeOut();
        }
    });
    var backToTop = $('#back-to-top')
    // scroll body to 0px on click
    backToTop.click(function () {
        backToTop.tooltip('hide');
        $('body,html').animate({
            scrollTop: 0
        }, 800);
        return false;
    });
    backToTop.tooltip('show');

    $(document).on("click", ".j-verifycode", function () {
        var path = $(this).attr("src");
        var index = path.indexOf("?");
        path = (index == -1) ? (path + "?" + new Date()) : path.substring(0, index + 1) + new Date();
        $(this).attr("src", path);
    });
    var messageBox = $('#messagebox');
    $('#commentForm').ajaxForm(function (data) {
        if (data.succeed) {
            window.location.reload();
        } else {
            messageBox.show();
            setTimeout(hideMessagebox, 2000);
            messageBox.html(data.message);
            $("input[name='verifyCode']").val('');
        }
        function hideMessagebox() {
            messageBox.hide();
        }
    });

    // $("#articleDelete").click(function (event) {
    //     if (confirm("Are you sure to delete?")) {
    //         articleDelete($("#articleId").text());
    //     }
    // });
});