package frontend

const indexTemplate string = `
<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="utf-8">
        <meta http-equiv="X-UA-Compatible" content="IE=edge">
        <meta name="viewport" content="width=device-width, initial-scale=1">
        <title>DDNS</title>

        <!-- Latest compiled and minified CSS -->
        <link rel="stylesheet" href="//maxcdn.bootstrapcdn.com/bootstrap/3.2.0/css/bootstrap.min.css">
        <link rel="stylesheet" href="//maxcdn.bootstrapcdn.com/font-awesome/4.1.0/css/font-awesome.min.css">

        <!-- Optional theme -->
        <style type="text/css" media="all">
            /* Space out content a bit */
            body {
                padding-top: 20px;
                padding-bottom: 20px;
            }

            /* Everything but the jumbotron gets side spacing for mobile first views */
            .header,
            .marketing,
            .footer {
                padding-right: 15px;
                padding-left: 15px;
            }

            /* Custom page header */
            .header {
                border-bottom: 1px solid #e5e5e5;
            }
            /* Make the masthead heading the same height as the navigation */
            .header h3 {
                padding-bottom: 19px;
                margin-top: 0;
                margin-bottom: 0;
                line-height: 40px;
            }

            /* Custom page footer */
            .footer {
                padding-top: 19px;
                color: #777;
                border-top: 1px solid #e5e5e5;
            }

            /* Customize container */
            @media (min-width: 768px) {
                .container {
                    max-width: 730px;
                }
            }
            .container-narrow > hr {
                margin: 30px 0;
            }

            /* Main marketing message and sign up button */
            .jumbotron {
                text-align: center;
                border-bottom: 1px solid #e5e5e5;
            }
            .jumbotron .btn {
                padding: 14px 24px;
                font-size: 21px;
            }

            /* Supporting marketing content */
            .marketing {
                margin: 40px 0;
            }
            .marketing p + h4 {
                margin-top: 28px;
            }

            /* Responsive: Portrait tablets and up */
            @media screen and (min-width: 768px) {
                /* Remove the padding we set earlier */
                .header,
                .marketing,
                .footer {
                    padding-right: 0;
                    padding-left: 0;
                }
                /* Space out the masthead */
                .header {
                    margin-bottom: 30px;
                }
                /* Remove the bottom border on the jumbotron for visual effect */
                .jumbotron {
                    border-bottom: 0;
                }
            }

        </style>

        <!-- HTML5 Shim and Respond.js IE8 support of HTML5 elements and media queries -->
        <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
        <!--[if lt IE 9]>
        <script src="https://oss.maxcdn.com/html5shiv/3.7.2/html5shiv.min.js"></script>
        <script src="https://oss.maxcdn.com/respond/1.4.2/respond.min.js"></script>
        <![endif]-->
    </head>
    <body>
        <div class="container">
            <div class="header">
                <ul class="nav nav-pills pull-right">
                    <li class="active"><a href="#">Home</a></li>
                    <li><a href="https://github.com/pboehm/ddns">
                        <i class="fa fa-github fa-lg"></i> Code</a></li>
                </ul>
                <h3 class="text-muted">DDNS</h3>
            </div>

            <div class="jumbotron">
                <h1>Self-Hosted Dynamic DNS</h1>

                <p class="lead">DDNS is a project that lets you host a Dynamic
                DNS Service, similar to DynDNS/NO-IP, on your own servers.</p>

                <hr />

                <form class="form-inline" role="form">
                    <div id="hostname_input" class="form-group">
                        <div class="input-group">
                            <input id="hostname" class="form-control input-lg" type="text" placeholder="my-own-hostname">
                            <div class="input-group-addon input-lg">{{.domain}}</div>
                        </div>
                    </div>
                </form>

                <hr />

                <input type="button" id="register" class="btn btn-primary disabled" value="Register Host" />
            </div>

            <div id="command_output"></div>

            <div class="footer">
                <p>&copy; Philipp Böhm</p>
            </div>

        </div> <!-- /container -->

        <!-- jQuery (necessary for Bootstrap's JavaScript plugins) -->

        <script src="https://ajax.googleapis.com/ajax/libs/jquery/1.11.1/jquery.min.js"></script>
        <!-- Latest compiled and minified JavaScript -->
        <script src="//maxcdn.bootstrapcdn.com/bootstrap/3.2.0/js/bootstrap.min.js"></script>

        <script type="text/javascript" charset="utf-8">

            function isValid() {
                $('#register').removeClass("disabled");
                $('#hostname_input').removeClass("has-error");
                $('#hostname_input').addClass("has-success");
            }

            function isNotValid(argument) {
                $('#register').addClass("disabled");
                $('#hostname_input').removeClass("has-success");
                $('#hostname_input').addClass("has-error");
            }

            function validate() {
                var hostname = $('#hostname').val();

                $.getJSON("/available/" + hostname, function( data ) {
                    if (data.available) {
                        isValid();
                    } else {
                        isNotValid();
                    }
                }).error(function(){ isNotValid(); });
            }

            $(document).ready(function() {
                var timer = null;
                $('#hostname').on('keydown', function () {
                    clearTimeout(timer);
                    timer = setTimeout(validate, 800)
                });


                $('#register').click(function() {
                    var hostname = $("#hostname").val();

                    $.getJSON("/new/" + hostname, function( data ) {
                        console.log(data);

                        var host = location.protocol + '//' + location.host;

                        $("#command_output").append(
                            "<pre>curl \"" + host +
                            data.update_link + "\"</pre>");
                    })
                });
            });
        </script>
    </body>
</html>
`
