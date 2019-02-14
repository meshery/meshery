// Example starter JavaScript for disabling form submissions if there are invalid fields
    // source: https://getbootstrap.com/docs/4.1/examples/checkout/?
var customEditor;
(function() {
    'use strict';

    window.addEventListener('load', function() {
        // Fetch all the forms we want to apply custom Bootstrap validation styles to
        var forms = document.getElementsByClassName('needs-validation');

        // Loop over them and prevent submission
        var validation = Array.prototype.filter.call(forms, function(form) {
        form.addEventListener('submit', function(event) {
            if (form.checkValidity() === false) {
            event.preventDefault();
            event.stopPropagation();
            }
            form.classList.add('was-validated');
        }, false);
        });
    }, false);

    $('#loadTestSuccess').hide();
    $('#loadTestError').hide();

    $('#istioConfigSuccess').hide();
    $('#istioConfigError').hide();


    // just to render the charts initially
    var data = fortioResultToJsChartData({
        DurationHistogram: {
            Count: 0,
            Data: [
                {
                    Start: new Date(),
                }
            ],
            Max: 0,
            Min: 0,
            Avg: 0,
            Percentiles: []
        },
        StartTime: new Date(),
        URL: '',
        Labels: '',
        RetCodes: {}
    });
    showChart(data);

    $("#loadTestForm").on( "submit", function( event ) {
        event.preventDefault();
        $('#loadRun').addClass("disabled");
        var loadTestForm = $(this)[0];
        if ((typeof(loadTestForm.checkValidity) == "function") && !loadTestForm.checkValidity()) {
            return;
        }
        startCountDownCounter(loadTestForm.t.value);

        var data = $( this ).serialize();
        $("#loadTestSuccess").hide();
        $("#loadTestError").hide();
        console.log("serialized data: "+data);
        jQuery.ajax({
            url: "/play/load-test",
            type: "POST", 
            data: data, 
            error: function(){
                $('#loadRun').removeClass('disabled');
                $("#loadTestError").fadeTo(2000, 500); 
                resetCounter();
            },
            success: function(result){
                $('#loadRun').removeClass('disabled');
                $("#loadTestSuccess").fadeTo(2000, 500).slideUp(500, function(){
                    $("#loadTestSuccess").slideUp(500);
                    $('#ldClose').click();
                }); 
                var res = JSON.parse(result);
                $('#loadTestRunTime').text(new Date(res.StartTime));
                var data = fortioResultToJsChartData(res);
                showChart(data);
                // $(".loadTestResults").show();
            }
        });
    });

    $("#istioForm").on('submit', function( event ) {
        event.preventDefault();
        var istioForm = $(this)[0];
        if ((typeof(istioForm.checkValidity) == "function") && !istioForm.checkValidity()) {
            return;
        }
        $('#customBody').val(customEditor.getValue());
        var data = $(this).serialize();
        // console.log("serialized data: "+data);
        $("#istioConfigError").hide();
        $("#istioConfigSuccess").hide();
        jQuery.ajax({
            url: "/play/mesh",
            type: "POST", 
            data: data, 
            error: function(resp){
                if (resp && resp.responseText) {
                    $('#istioConfigErrorResp').html(resp.responseText);
                }
                $("#istioConfigError").fadeTo(2000, 500); 
            },
            success: function(result){
                $("#istioConfigSuccess").fadeTo(2000, 500).slideUp(500, function(){
                    $("#istioConfigSuccess").slideUp(500);
                    $('#isClose').click();
                }); 
            }
        });
    });

    customEditor = CodeMirror.fromTextArea(document.getElementById("customBody"), {
        lineNumbers: true,
        lineWrapping: true,
        gutters: ["CodeMirror-lint-markers"],
        lint: true,
        theme: "monokai",
        mode: "text/x-yaml"
    });

    $('input[type="radio"][name="query"]').change(function(){
        if ($(this).is(':checked') && $(this).val() == "custom") {
            $('#customBodyDiv').show();
        } else {
            $('#customBodyDiv').hide();
        }
    });


})();



function resetCounter(){
    $('#countdown #hour').text('00');
    $('#countdown #mins').text('00');
    $('#countdown #secs').text('00');
}

function startCountDownCounter(startMin){
    resetCounter();
    if(startMin.length == 1){
        startMin = '0' + startMin;
    }
    $('#ltCountDown #mins').html(startMin);
    var startTime = new Date();

    var intv = setInterval(function time(){
        var d = startMin * 60 - Math.abs(new Date().getTime() - startTime.getTime()) / 1000;

        var days = Math.floor(d / 86400);
        d -= days * 86400;

        var hrs = Math.floor(d / 3600) % 24;
        d -= hrs * 3600;

        var mins = Math.floor(d / 60) % 60;
        d -= mins * 60;

        var secs = Math.floor(d % 60);

        var hr = hrs + '';
        if(hr.length == 1){
            hr = '0' + hrs;
        }

        var min = mins + '';
        if(min.length == 1){
            min = '0' + mins;
        }
        var sec = secs + '';
        if(sec.length == 1){
            sec = '0' + secs;
        }
        $('#ltCountDown #hour').html(hr);
        $('#ltCountDown #mins').html(min);
        $('#ltCountDown #secs').html(sec);

        if ($('#ltCountDown #mins').html() == "00" && $('#ltCountDown #secs').html() == "00") {
        clearInterval(intv);
        }

    }, 1000);
}