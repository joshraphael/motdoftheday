$(document).ready(function () {
    $(".history").hide()
    var init_id = $("#history").children(":selected").attr("id")
    $("#history-" + init_id).show()
    $("#history").change(function () {
        var id = $(this).children(":selected").attr("id")
        $(".history").hide();
        $("#history-" + id).show()
    });
})