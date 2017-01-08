/* global clear_blocks */
/* global formatMoney */
/* global in_array */
/* global new_block */
/* global formatDate */
/* global nDig */
/* global randStr */
/* global bag */
/* global $ */
var ws = {};
var user = {
    username: bag.session.username,
    name: bag.session.name,
    role: bag.session.role
};

var panels = [
    {
        name: "title",
        formID: "tradeFilter",
        tableID: "#tradesBody",
        filterPrefix: "trade_"
    },
    {
        name: "license",
        formID: "licenseFilter",
        tableID: "#licensesBody",
        filterPrefix: "license_"
    },
    {
        name: "registration",
        formID: "registrationFilter",
        tableID: "#registrationBody",
        filterPrefix: "registration_"
    },
    {
        name: "toll",
        formID: "tollFilter",
        tableID: "#tollBody",
        filterPrefix: "toll_"
    },
    {
        name: "violation",
        formID: "violationFilter",
        tableID: "#violationBody",
        filterPrefix: "violation_"
    },
    {
        name: "audit",
        formID: "auditFilter",
        tableID: "#auditBody",
        filterPrefix: "audit_"
    },
    {
        name: "register-vehicle",
        formID: "register-vehicleFilter",
        tableID: "#register-vehicleBody",
        filterPrefix: "register-vehicle_"
    },
    {
        name: "issue-driver-license",
        formID: "issue-driver-licenseFilter",
        tableID: "#issue-driver-licenseBody",
        filterPrefix: "issue-driver-license_"
    },
    {
        name: "issue-title",
        formID: "issue-titleFilter",
        tableID: "#issue-titleBody",
        filterPrefix: "issue-title_"
    }
];

// =================================================================================
// On Load
// =================================================================================
$(document).on('ready', function () {
    connect_to_server();
    if (user.name) $("#userField").html(user.name.toUpperCase() + ' ');

    // Customize which panels show up for which user
    $(".nav").hide();
    console.log("user role", bag.session.user_role);

    // Only show tabs if a user is logged in
    if (user.username) {

        // Display tabs based on user's role
        if (user.role && user.role.toUpperCase() === "auditor".toUpperCase()) {
            $("#auditLink").show();
        } else if (user.username) {
            //$("#issue-titleLink").show();
            //$("#issue-driver-licenseLink").show();
            //$("#register-vehicleLink").show();
            $("#titleLink").show();
            $("#licenseLink").show();
            $("#registrationLink").show();
            $("#tollLink").show();
            $("#traffic-violationLink").show();
        }
    } else {

        // Display the login and user registration links
        $("#loginLink").show();
        $("#registerLink").show();
    }

    // =================================================================================
    // jQuery UI Events
    // =================================================================================

    $("#submit1").click(function () {
         var myDate = new Date();
         myDate.setFullYear(myDate.getFullYear() + 2);
         var obj = {
                type: "create",
                registration: {
                    registrationId: Date.now().toString(),
                    plateNum: Date.now().toString(),
                    vin: escapeHtml($("input[name='vehicleId']").val()),
                    testId: escapeHtml($("input[name='inspectionReceiptId']").val()),
                    policyId: escapeHtml($("input[name='insurancePolicyId']").val()),
                    owner: user.name,
                    auto: escapeHtml($("select[name='auto-renew-r']").val()),
                    issueDate: Date.now().toString(),
                    expiryDate: myDate.toString()
                },
                user: user.username
            };
            if (obj.registration && obj.registration.testId) {
                console.log('creating registration, sending', obj);
                //alert(JSON.stringify(obj));
                ws.send(JSON.stringify(obj));
                $(".panel").hide();
                $("#registrationPanel").show();
            }
   });

   $("#submitIssueDriverLicense").click(function () {
         var myDate = new Date(); 
         myDate.setFullYear(myDate.getFullYear() + 2);
         var obj = {
                type: "create",
                license: {
                    licenseId: Date.now().toString(),
                    testId: escapeHtml($("input[name='roadTestReceiptId']").val()),
                    address: escapeHtml($("input[name='address']").val()),
                    city: escapeHtml($("input[name='city']").val()),
                    state: escapeHtml($("input[name='state']").val()),
                    zip: escapeHtml($("input[name='zip']").val()),
                    driver: user.name,
                    auto: escapeHtml($("select[name='auto-renew']").val()),
                    issueDate: Date.now().toString(),
                    expiryDate: myDate.toString()
                },
                user: user.username
            };
            if (obj.license && obj.license.testId) {
                console.log('creating license, sending', obj);
                //alert(JSON.stringify(obj));
                ws.send(JSON.stringify(obj));
                $(".panel").hide();
                $("#licensePanel").show();
            }
   });

   $("#submit").click(function () {
        if (user.username) {
	    
            var obj = {
                type: "create",
                paper: {
                    vin: escapeHtml($("input[name='assetId']").val()),
                    make: escapeHtml($("select[name='make']").val()),
                    model: escapeHtml($("input[name='vehicle-model']").val()),
                    year: Number($("select[name='year']").val()),
                    color: escapeHtml($("select[name='vehicle-color']").val()),
                    miles: Number($("input[name='vehicle-miles']").val()),
                    value: Number($("input[name='vehicle-value']").val()),
                    owner: user.name,
                    issuer: user.name,
                    state: "active",
                    issueDate: Date.now().toString()
                }, 
                user: user.username
            };
            if (obj.paper && obj.paper.vin) {
                console.log('creating vehicle, sending', obj);
                //alert(JSON.stringify(obj));
                ws.send(JSON.stringify(obj));
                $(".panel").hide();
                $("#titlePanel").show();
	    }
          
        }
        return false;
    });

    $("#issue-titleLink").click(function () {
        $("input[name='name']").val('r' + randStr(6));
    });

    $("#issue-driver-licenseLink").click(function () {
        $("input[name='name']").val('r' + randStr(6));
    });
    $("#register-vehicleLink").click(function () {
        $("input[name='name']").val('r' + randStr(6));
    });
    $("#createLink").click(function () {
        $("input[name='name']").val('r' + randStr(6));
    });


    $("#titleLink").click(function () {
        ws.send(JSON.stringify({type: "get_open_trades", v: 2, user: user.username}));
    });
    $("#licenseLink").click(function () {
        ws.send(JSON.stringify({type: "get_open_licenses", v: 2, user: user.username}));
    });
    $("#registrationLink").click(function () {
        ws.send(JSON.stringify({type: "get_open_registrations", v: 2, user: user.username}));
    });
    $("#tollLink").click(function () {
        ws.send(JSON.stringify({type: "get_open_tolls", v: 2, user: user.username}));
    });

    $("#traffic-violationLink").click(function () {
        ws.send(JSON.stringify({type: "get_open_violations", v: 2, user: user.username}));
    });

    //login events
    $("#whoAmI").click(function () {				//drop down for login
        if ($("#loginWrap").is(":visible")) {
            $("#loginWrap").fadeOut();
        }
        else {
            $("#loginWrap").fadeIn();
        }
    });

    // Filter the trades whenever the filter modal changes
    $(".trade-filter").keyup(function () {
        "use strict";
        console.log("Change in trade filter detected.");
        processFilterForm(panels[0]);
    });
    $(".audit-filter").keyup(function () {
        "use strict";
        console.log("Change in audit filter detected.");
        processFilterForm(panels[1]);
    });

    // Click events for the columns of the table
    $('.sort-selector').click(function () {
        "use strict";
        var sort = $(this).attr('sort');

        // Clear any sort direction arrows
        $('span').remove('.sort-indicator');

        // Clicking the column again should reverse the sort
        if (sort_papers[sort] === sort_selected) {
            console.log("Reversing the table");
            sort_reversed = !sort_reversed;
        }
        else sort_reversed = false;

        // Add the appropriate arrow to the current selector
        var arrow_icon = (sort_reversed ? 'fa-arrow-up' : 'fa-arrow-down');
        var span = document.createElement('span');
        span.classList.add('fa');
        span.classList.add(arrow_icon);
        span.classList.add('sort-indicator');
        $(this).append(span);

        // Change to the sort corresponding to that column
        sort_selected = sort_papers[sort];
        console.log("Sorting by:", sort);
        for (var i in panels) {
            build_trades(bag.papers, panels[i]);
        }
    });

    //trade events
    $(document).on("click", ".buyPaper", function () {
        if (user.username) {
            console.log('transfer title...');
            var vin = $(this).attr('data_vin');
            var issuer = $(this).attr('data_issuer');
            var newOwner = escapeHtml($("input[name='newOwner']").val());
            var amount = Number($("input[name='amtPaid']").val());

            var msg = {
                type: 'transfer_paper',
                transfer: {
                    vin: vin,
                    fromOwner: issuer,
                    toOwner: newOwner,
                    issueDate: Date.now().toString(),
                    amountPaid: amount
                },
                user: user.username
            };
            console.log('sending transfer title', msg);
	    //alert(JSON.stringify(msg));
            ws.send(JSON.stringify(msg));
            $("#notificationPanel").animate({width: 'toggle'});
        }
    });

    //terminate asset
    $(document).on("click", ".terminateAsset", function () {
        if (user.username) {
            console.log('terminating asset...');
            var vin = $(this).attr('data_vin');
            var issuer = $(this).attr('data_issuer');

            var msg = {
                type: 'terminate_asset',
                terminate: {
                    vin: vin,
                    owner: issuer,
                    issueDate: Date.now().toString()
                },
                user: user.username
            };
            console.log('sending terminate asset', msg);
            //alert(JSON.stringify(msg));
            ws.send(JSON.stringify(msg));
            $("#vehicleTerminationNotificationPanel").animate({width: 'toggle'});
        }
    });

   // send mail
   $(document).on("click", ".sendMail", function () {
        var sendTo = $(this).attr('data_owner') + "@gmail.com";
        var letter = "Dear " + $(this).attr('data_owner') + " This ticket " + $(this).attr('data_id') + " is with regards to a traffic violation namely " + $(this).attr('data_type') + " that you committed on " + $(this).attr('data_date') + " @ location " + $(this).attr('data_loc') + ". A fine amount of AED " + $(this).attr('data_amt') + " has been charged to your account"
        var msg = {
                type: 'send_mail',
                mail: {
		    to: sendTo,
		    text: letter
                },
                user: user.username
            };
            console.log('sending', msg);
            ws.send(JSON.stringify(msg));
   });

   //renew license
    $(document).on("click", ".renewLicense", function () {
        var myDate = new Date(); 
        myDate.setFullYear(myDate.getFullYear() + 2);
        if (user.username) {
            console.log('renewing license...');
            //var i = $(this).attr('trade_pos');
            var licenseId = $(this).attr('data_licenseId');
            var driver = $(this).attr('data_driver');


            // TODO Map the trade_pos to the correct button
            var msg = {
                type: 'renew_license',
                renewlicense: {
                    txId: Date.now().toString(),
                    licenseId: licenseId,
                    driver: driver,
                    issueDate: Date.now().toString(),
                    expiryDate: myDate.toString()
                },
                user: user.username
            };
            console.log('sending', msg);
            ws.send(JSON.stringify(msg));
            $("#licenseRenewalNotificationPanel").animate({width: 'toggle'});
        }
    });

    //simulate traffic violation
    $(document).on("click", ".simulateViolation", function () {
         var myDate = new Date();
         myDate.setFullYear(myDate.getFullYear() + 2);
        if (user.username) {
            console.log('simulating traffic violation...');
            var licenseId = $(this).attr('data_licenseId');
            var driver = $(this).attr('data_driver');
            var amt = 60;
	    var loc = "3 N 50 S";
	    var type = "speeding";

            var msg = {
                type: 'simulate_violation',
                trafficViolation: {
                    txId: Date.now().toString(),
		    type: type,
                    licenseId: licenseId,
                    driver: driver,
                    issueDate: Date.now().toString(),
                    fine: amt,
		    location: loc
                },
                user: user.username
            };
            console.log('sending', msg);
            ws.send(JSON.stringify(msg));
            $("#trafficViolationNotificationPanel").animate({width: 'toggle'});
        }
    });

    //simulate notify TODO
    $(document).on("click", ".simulateNotify", function () {
        if (user.username) {
            console.log('simulating notify...');
            var licenseId = $(this).attr('data_licenseId');
            var driver = $(this).attr('data_driver');
            var amt = 60;
            var loc = "3 N 50 S";
            var type = "speeding";

            var msg = {
                type: 'simulate_violation',
                notify: {
                    txId: Date.now().toString(),
                    type: type,
                    licenseId: licenseId,
                    driver: driver,
                    issueDate: Date.now().toString(),
                    fine: amt,
                    location: loc
                },
                user: user.username
            };
            console.log('sending', msg);
            ws.send(JSON.stringify(msg));
            $("#notifyNotificationPanel").animate({width: 'toggle'});
        }
    });

    //renew registration
    $(document).on("click", ".renewRegistration", function () {
         var myDate = new Date();
         myDate.setFullYear(myDate.getFullYear() + 4);
        if (user.username) {
            console.log('renewing registration...');
            var registrationId = $(this).attr('data_registrationId');
            var owner = $(this).attr('data_owner');

            var msg = {
                type: 'renew_registration',
                renewregistration: {
                    txId: Date.now().toString(),
                    registrationId: registrationId,
                    owner: owner,
                    issueDate: Date.now().toString(),
                    expiryDate: myDate.toString()
                },
                user: user.username
            };
            console.log('sending', msg);
            ws.send(JSON.stringify(msg));
            $("#registrationRenewalNotificationPanel").animate({width: 'toggle'});
        }
    });

    //simulate toll
    $(document).on("click", ".simulateToll", function () {
         var myDate = new Date();
         myDate.setFullYear(myDate.getFullYear() + 4);
        if (user.username) {
            console.log('renewing registration...');
            var registrationId = $(this).attr('data_registrationId');
            var owner = $(this).attr('data_owner');
	    var amt = 5;
            var loc = "3 N 50 S";
            var type = "highway toll";

            var msg = {
                type: 'simulate_toll',
                toll: {
                    txId: Date.now().toString(),
                    type: type,
                    registrationId: registrationId,
		    owner: owner,
                    issueDate: Date.now().toString(),
		    tollAmt: amt,
		    location: loc
                },
                user: user.username
            };
            console.log('sending', msg);
            ws.send(JSON.stringify(msg));
            $("#tollNotificationPanel").animate({width: 'toggle'});
        }
    });
});



// =================================================================================
// Helper Fun
// =================================================================================
function escapeHtml(str) {
    var div = document.createElement('div');
    div.appendChild(document.createTextNode(str));
    return div.innerHTML;
}

// =================================================================================
// Socket Stuff
// =================================================================================
function connect_to_server() {
    var connected = false;
    connect();

    function connect() {
        var wsUri = '';
        console.log('protocol', window.location.protocol);
        if (window.location.protocol === 'https:') {
            wsUri = "wss://" + bag.setup.SERVER.EXTURI;
        }
        else {
            wsUri = "ws://" + bag.setup.SERVER.EXTURI;
        }

        ws = new WebSocket(wsUri);
        ws.onopen = function (evt) {
            onOpen(evt);
        };
        ws.onclose = function (evt) {
            onClose(evt);
        };
        ws.onmessage = function (evt) {
            onMessage(evt);
        };
        ws.onerror = function (evt) {
            onError(evt);
        };
    }

    function onOpen(evt) {
        console.log("WS CONNECTED");
        connected = true;
        clear_blocks();
        $("#errorNotificationPanel").fadeOut();
        ws.send(JSON.stringify({type: "chainstats", v: 2, user: user.username}));
        ws.send(JSON.stringify({type: "get_papers", v: 2, user: user.username}));
        if (user.name && user.role !== "auditor") {
            ws.send(JSON.stringify({type: 'get_company', company: user.name, user: user.username}));
        }
    }

    function onClose(evt) {
        console.log("WS DISCONNECTED", evt);
        connected = false;
        setTimeout(function () {
            connect();
        }, 5000);					//try again one more time, server restarts are quick
    }

	function onMessage(msg) {
		try {
			var data = JSON.parse(msg.data);
			console.log('rec', data);
			if (data.msg === 'papers') {
				try{
					var papers = JSON.parse(data.papers);
					build_trades(papers, panels[0]);
				}
				catch(e){
					console.log('cannot parse papers', e);
				}
			}
			else if (data.msg === 'licenses') {
                                try{
                                        var licenses = JSON.parse(data.licenses);
                                        build_licenses(licenses, panels[1]);
                                }
                                catch(e){
                                        console.log('cannot parse licenses', e);
                                }
                        }
                        else if (data.msg === 'registrations') {
                                try{
                                        var registrations = JSON.parse(data.registrations);
                                        build_registrations(registrations, panels[2]);
                                }
                                catch(e){
                                        console.log('cannot parse licenses', e);
                                }
                        }
                        else if (data.msg === 'tolls') {
                                try{
                                        var tolls = JSON.parse(data.tolls);
                                        build_tolls(tolls, panels[3]);
                                }
                                catch(e){
                                        console.log('cannot parse tolls', e);
                                }
                        }
                        else if (data.msg === 'violations') {
                                try{
                                        var violations = JSON.parse(data.violations);
                                        build_violations(violations, panels[4]);
                                }
                                catch(e){
                                        console.log('cannot parse violations', e);
                                }       
                        } 
			else if (data.msg === 'chainstats') {
				console.log(JSON.stringify(data));
				var e = formatDate(data.blockstats.transactions[0].timestamp.seconds * 1000, '%M/%d/%Y &nbsp;%I:%m%P');
				$("#blockdate").html('<span style="color:#fff">TIME</span>&nbsp;&nbsp;' + e + ' UTC');
				var temp = {
					id: data.blockstats.height,
					blockstats: data.blockstats
				};
				new_block(temp);									//send to blockchain.js
			}
			else if (data.msg === 'company') {
				try{
					var company = JSON.parse(data.company);
					$("#accountBalance").html(formatMoney(company.cashBalance));
					$("#accountBalanceLicense").html(formatMoney(company.cashBalance));
					$("#accountBalanceRegistration").html(formatMoney(company.cashBalance));
					$("#accountBalanceToll").html(formatMoney(company.cashBalance));
					$("#accountBalanceViolation").html(formatMoney(company.cashBalance));
				}
				catch(e){
					console.log('cannot parse company', e);
				}
			}
			else if (data.msg === 'reset') {
				// Ask for all available trades and information for the current company
				ws.send(JSON.stringify({type: "get_papers", v: 2, user: user.username}));
				ws.send(JSON.stringify({type: "get_licenses", v: 2, user: user.username}));
				ws.send(JSON.stringify({type: "get_registrations", v: 2, user: user.username}));
				ws.send(JSON.stringify({type: "get_tolls", v: 2, user: user.username}));
				ws.send(JSON.stringify({type: "get_violations", v: 2, user: user.username}));
                		ws.send(JSON.stringify({type: "chainstats", v: 2, user: user.username}));
				if (user.role !== "auditor") {
					ws.send(JSON.stringify({type: 'get_company', company: user.name, user: user.username}));
				}
			}
			else if (data.type === 'error') {
				console.log("Error:", data.error);
			}
		}
		catch (e) {
			console.log('ERROR', e);
			//ws.close();
		}
	}

    function onError(evt) {
        console.log('ERROR ', evt);
        if (!connected && bag.e == null) {		//don't overwrite an error message
            $("#errorName").html("Warning");
            $("#errorNoticeText").html("Waiting on the node server to open up so we can talk to the blockchain. ");
            $("#errorNoticeText").append("This app is likely still starting up. ");
            $("#errorNoticeText").append("Check the server logs if this message does not go away in 1 minute. ");
            $("#errorNotificationPanel").fadeIn();
        }
    }

    function sendMessage(message) {
        console.log("SENT: " + message);
        ws.send(message);
    }
}


// =================================================================================
//	UI Building
// =================================================================================
function build_violations(violations, panelDesc) {

    if (violations && violations.length > 0) {

        // Break the violations down into entries
        console.log('breaking violations into individual entries');
        var entries = [];
        for (var violation in violations) {
            var broken_up = violation_to_entries(violations[violation]);
            entries = entries.concat(broken_up);
        }
        console.log("Displaying", violations.length, "violation as", entries.length, "entries");

        // If no panel is given, assume this is the 
        if (!panelDesc) {
            panelDesc = panels[0];
        }

        entries.sort(sort_selected);
        if (sort_reversed) entries.reverse();

        // Display each entry as a row in the table
        var rows = [];
        for (var i in entries) {
            console.log('!', entries[i]);


                if (excluded(entries[i], filter)) {
                    var style;

                    if (user.name.toLowerCase() === entries[i].owner.toLowerCase() || user.name.toLowerCase() === 'government') {
                        style = null;
                    }else {
                        style = 'invalid';
                    }

                    // Create a row for each valid trade
                    var data = [
                        formatDate(Number(entries[i].issueDate), '%M/%d %I:%m%P'),
                        entries[i].txId,
                        entries[i].type,
                        entries[i].licenseId,
                        escapeHtml(entries[i].owner),
                        entries[i].tollAmt,
			entries[i].location
                    ];

                    var row = createRow(data);
                    style && row.classList.add(style);

                    if (panelDesc.name === "violation") {

                        var simulate = true;
                        if (user.name.toLowerCase() === 'government') simulate = false;
                        var button1 = simulateNotifyButton(simulate, entries[i]);
                        row.appendChild(button1);
                    }
                    rows.push(row);
                }

        }

        // Placeholder for an empty table
        var html = '';
        if (rows.length == 0) {
            if (panelDesc.name === 'violation')
                html = '<tr><td>nothing here...</td><td></td><td></td><td></td><td></td><td></td><td></td><td></td><td></td><td></td></tr>';
            else if (panelDesc.name === 'audit')
                html = '<tr><td>nothing here...</td><td></td><td></td><td></td><td></td><td></td><td></td><td></td><td></td></tr>'; // No action column
            $(panelDesc.tableID).html(html);
        } else {
            // Remove the existing table data
            console.log("clearing existing table data");
            var tableBody = $(panelDesc.tableID);
            tableBody.empty();


            // Add the new rows to the table
            console.log("populating new table data");
            var row;
            while (rows.length > 0) {
                row = rows.shift();
                tableBody.append(row);
            }
        }
    } else {
        if (panelDesc.name === 'violation')
            html = '<tr><td>nothing here...</td><td></td><td></td><td></td><td></td><td></td><td></td><td></td><td></td><td></td></tr>';
        else if (panelDesc.name === 'audit')
            html = '<tr><td>nothing here...</td><td></td><td></td><td></td><td></td><td></td><td></td><td></td><td></td></tr>'; // No action column
        $(panelDesc.tableID).html(html);
    }
}

function build_tolls(tolls, panelDesc) {

    if (tolls && tolls.length > 0) {

        // Break the tolls down into entries
        console.log('breaking tolls into individual entries');
        var entries = [];
        for (var toll in tolls) {
            var broken_up = toll_to_entries(tolls[toll]);
            entries = entries.concat(broken_up);
        }
        console.log("Displaying", tolls.length, "toll as", entries.length, "entries");

        // If no panel is given, assume this is the 
        if (!panelDesc) {
            panelDesc = panels[0];
        }

        entries.sort(sort_selected);
        if (sort_reversed) entries.reverse();

        // Display each entry as a row in the table
        var rows = [];
        for (var i in entries) {
            console.log('!', entries[i]);


                if (excluded(entries[i], filter)) {
                    var style;

                    if (user.name.toLowerCase() === entries[i].owner.toLowerCase() || user.name.toLowerCase() === 'government') {
                        style = null;
                    }else {
                        style = 'invalid';
                    }

                    // Create a row for each valid trade
                    var data = [
                        formatDate(Number(entries[i].issueDate), '%M/%d %I:%m%P'),
                        entries[i].txId,
                        entries[i].type,
                        entries[i].plateNum,
                        escapeHtml(entries[i].owner),
                        entries[i].tollAmt,
			entries[i].location
                    ];

                    var row = createRow(data);
                    style && row.classList.add(style);

                    if (panelDesc.name === "toll") {

                        var simulate = true;
                        if (user.name.toLowerCase() === 'government') simulate = false;
                        var button1 = simulateNotifyButton(simulate, entries[i].registrationId, entries[i].owner);
                        row.appendChild(button1);
                    }
                    rows.push(row);
                }

        }

        // Placeholder for an empty table
        var html = '';
        if (rows.length == 0) {
            if (panelDesc.name === 'toll')
                html = '<tr><td>nothing here...</td><td></td><td></td><td></td><td></td><td></td><td></td><td></td><td></td><td></td></tr>';
            else if (panelDesc.name === 'audit')
                html = '<tr><td>nothing here...</td><td></td><td></td><td></td><td></td><td></td><td></td><td></td><td></td></tr>'; // No action column
            $(panelDesc.tableID).html(html);
        } else {
            // Remove the existing table data
            console.log("clearing existing table data");
            var tableBody = $(panelDesc.tableID);
            tableBody.empty();


            // Add the new rows to the table
            console.log("populating new table data");
            var row;
            while (rows.length > 0) {
                row = rows.shift();
                tableBody.append(row);
            }
        }
    } else {
        if (panelDesc.name === 'toll')
            html = '<tr><td>nothing here...</td><td></td><td></td><td></td><td></td><td></td><td></td><td></td><td></td><td></td></tr>';
        else if (panelDesc.name === 'audit')
            html = '<tr><td>nothing here...</td><td></td><td></td><td></td><td></td><td></td><td></td><td></td><td></td></tr>'; // No action column
        $(panelDesc.tableID).html(html);
    }
}

function build_registrations(registrations, panelDesc) {

    if (registrations && registrations.length > 0) {

        // Break the registrations down into entries
        console.log('breaking registrations into individual entries');
        var entries = [];
        for (var registration in registrations) {
            var broken_up = registration_to_entries(registrations[registration]);
            entries = entries.concat(broken_up);
        }
        console.log("Displaying", registrations.length, "registrations as", entries.length, "entries");

        // If no panel is given, assume this is the 
        if (!panelDesc) {
            panelDesc = panels[0];
        }

        entries.sort(sort_selected);
        if (sort_reversed) entries.reverse();

        // Display each entry as a row in the table
        var rows = [];
        for (var i in entries) {
            console.log('!', entries[i]);


                if (excluded(entries[i], filter)) {
                    var style;

                    if (user.name.toLowerCase() === entries[i].owner.toLowerCase() || user.name.toLowerCase() === 'government') {
                        style = null;
                    }else {
                        style = 'invalid';
                    }

                    // Create a row for each valid trade
                    var data = [
                        formatDate(Number(entries[i].issueDate), '%M/%d %I:%m%P'),
                        entries[i].registrationId,
                        entries[i].vin,
                        entries[i].plateNum,
                        entries[i].policyId,
                        escapeHtml(entries[i].owner),
                        entries[i].expiryDate,
                        entries[i].auto
                    ];

                    var row = createRow(data);
                    style && row.classList.add(style);

                    if (panelDesc.name === "registration") {
                        var disabled = false;
                        //if (user.name.toLowerCase() === entries[i].owner.toLowerCase()) disabled = false;		
                        var button = renewRegistrationButton(disabled, entries[i].registrationId, entries[i].owner);
                        row.appendChild(button);

                        var simulate = true;
                        if (user.name.toLowerCase() === 'government') simulate = false;
                        var button1 = simulateTollButton(simulate, entries[i].registrationId, entries[i].owner);
                        row.appendChild(button1);
                    }
                    rows.push(row);
                }

        }

        // Placeholder for an empty table
        var html = '';
        if (rows.length == 0) {
            if (panelDesc.name === 'registration')
                html = '<tr><td>nothing here...</td><td></td><td></td><td></td><td></td><td></td><td></td><td></td><td></td><td></td></tr>';
            else if (panelDesc.name === 'audit')
                html = '<tr><td>nothing here...</td><td></td><td></td><td></td><td></td><td></td><td></td><td></td><td></td></tr>'; // No action column
            $(panelDesc.tableID).html(html);
        } else {
            // Remove the existing table data
            console.log("clearing existing table data");
            var tableBody = $(panelDesc.tableID);
            tableBody.empty();


            // Add the new rows to the table
            console.log("populating new table data");
            var row;
            while (rows.length > 0) {
                row = rows.shift();
                tableBody.append(row);
            }
        }
    } else {
        if (panelDesc.name === 'registration')
            html = '<tr><td>nothing here...</td><td></td><td></td><td></td><td></td><td></td><td></td><td></td><td></td><td></td></tr>';
        else if (panelDesc.name === 'audit')
            html = '<tr><td>nothing here...</td><td></td><td></td><td></td><td></td><td></td><td></td><td></td><td></td></tr>'; // No action column
        $(panelDesc.tableID).html(html);
    }
}

function build_licenses(licenses, panelDesc) {

    if (licenses && licenses.length > 0) {

        // Break the licenses down into entries
        console.log('breaking licenses into individual entries');
        var entries = [];
        for (var license in licenses) {
            var broken_up = license_to_entries(licenses[license]);
            entries = entries.concat(broken_up);
        }
        console.log("Displaying", licenses.length, "licenses as", entries.length, "entries");

        // If no panel is given, assume this is the licenses panel
        if (!panelDesc) {
            panelDesc = panels[0];
        }

        entries.sort(sort_selected);
        if (sort_reversed) entries.reverse();

        // Display each entry as a row in the table
        var rows = [];
        for (var i in entries) {
            console.log('!', entries[i]);


                if (excluded(entries[i], filter)) {
                    var style;
                    if (user.name.toLowerCase() === entries[i].driver.toLowerCase() || user.name.toLowerCase() === 'government') {
                        //cannot buy my own stuff
                        style = null;
                    }else {
                        style = 'invalid';
                    }

                    // Create a row for each valid trade
                    var data = [
                        formatDate(Number(entries[i].issueDate), '%M/%d %I:%m%P'),
                        entries[i].licenseId,
                        escapeHtml(entries[i].driver),
                        entries[i].address,
                        entries[i].city,
                        entries[i].state,
                        entries[i].zip,
                        entries[i].expiryDate,
                        entries[i].auto
                    ];

                    var row = createRow(data);
                    style && row.classList.add(style);

                    if (panelDesc.name === "license") {
                        var disabled = false;
                        if (user.name.toLowerCase() === entries[i].driver.toLowerCase()) disabled = false;		

                        var button = renewLicenseButton(disabled, entries[i].licenseId, entries[i].driver);
                        row.appendChild(button);

			var simulate = true;
                        if (user.name.toLowerCase() === 'government') simulate = false;
			var button1 = simulateViolationButton(simulate, entries[i].licenseId, entries[i].driver);
                        row.appendChild(button1);
                    }
                    rows.push(row);
                }

        }

        // Placeholder for an empty table
        var html = '';
        if (rows.length == 0) {
            if (panelDesc.name === 'license')
                html = '<tr><td>nothing here...</td><td></td><td></td><td></td><td></td><td></td><td></td><td></td><td></td><td></td></tr>';
            else if (panelDesc.name === 'audit')
                html = '<tr><td>nothing here...</td><td></td><td></td><td></td><td></td><td></td><td></td><td></td><td></td></tr>'; // No action column
            $(panelDesc.tableID).html(html);
        } else {
            // Remove the existing table data
            console.log("clearing existing table data");
            var tableBody = $(panelDesc.tableID);
            tableBody.empty();


            // Add the new rows to the table
            console.log("populating new table data");
            var row;
            while (rows.length > 0) {
                row = rows.shift();
                tableBody.append(row);
            }
        }
    } else {
        if (panelDesc.name === 'license')
            html = '<tr><td>nothing here...</td><td></td><td></td><td></td><td></td><td></td><td></td><td></td><td></td><td></td></tr>';
        else if (panelDesc.name === 'audit')
            html = '<tr><td>nothing here...</td><td></td><td></td><td></td><td></td><td></td><td></td><td></td><td></td></tr>'; // No action column
        $(panelDesc.tableID).html(html);
    }
}
/**
 * Process the list of trades from the server and displays them in the trade list.
 * This function builds the tables for multiple panels, so an object is needed to
 * identify which table it should be drawing to.
 * @param papers The list of trades to display.
 * @param panelDesc An object describing what panel the trades are being shown in.
 */
function build_trades(papers, panelDesc) {

    if (papers && papers.length > 0) {

        // Break the papers down into entries
        console.log('breaking papers into individual entries');
        var entries = [];
        for (var paper in papers) {
            var broken_up = paper_to_entries(papers[paper]);
            entries = entries.concat(broken_up);
        }
        console.log("Displaying", papers.length, "papers as", entries.length, "entries");

        // If no panel is given, assume this is the trade panel
        if (!panelDesc) {
            panelDesc = panels[0];
        }

        entries.sort(sort_selected);
        if (sort_reversed) entries.reverse();

        // Display each entry as a row in the table
        var rows = [];
        for (var i in entries) {
            console.log('!', entries[i]);


                if (excluded(entries[i], filter)) {
                    var style;
                    if (user.name.toLowerCase() === entries[i].owner.toLowerCase() || user.name.toLowerCase() === 'government') {
                        //cannot buy my own stuff
                        style = null;
                    } else {
                        style = 'invalid';
                    }

                    console.log("Vehicle data " + entries[i].toString());
                    // Create a row for each valid trade
                    var data = [
                        formatDate(Number(entries[i].issueDate), '%M/%d %I:%m%P'),
                        entries[i].vin,
                        entries[i].make,
                        entries[i].model,
                        entries[i].year,
                        entries[i].color,
                        entries[i].miles,
                        entries[i].value,
                        entries[i].issuer,
                        entries[i].owner,
			entries[i].state
                    ];

                    var row = createRow(data);
                    style && row.classList.add(style);

                    // Only the trade panel should allow you to interact with trades
                    if (panelDesc.name === "title") {
                        var disabled = false;
                        if (user.name.toLowerCase() === entries[i].owner.toLowerCase()) disabled = false;			//cannot buy my own stuff
                        //if (entries[i].issuer.toLowerCase() !== entries[i].owner.toLowerCase()) disabled = true;
                        var button = buyButton(disabled, entries[i].vin, entries[i].owner);
                        row.appendChild(button);
                        var button1 = terminateButton(disabled, entries[i].vin, entries[i].owner);
                        row.appendChild(button1);
                    }
                    rows.push(row);
                }
            

        }

        // Placeholder for an empty table
        var html = '';
        if (rows.length == 0) {
            if (panelDesc.name === 'trade')
                html = '<tr><td>nothing here...</td><td></td><td></td><td></td><td></td><td></td><td></td><td></td><td></td><td></td></tr>';
            else if (panelDesc.name === 'audit')
                html = '<tr><td>nothing here...</td><td></td><td></td><td></td><td></td><td></td><td></td><td></td><td></td></tr>'; // No action column
            $(panelDesc.tableID).html(html);
        } else {
            // Remove the existing table data
            console.log("clearing existing table data");
            var tableBody = $(panelDesc.tableID);
            tableBody.empty();


            // Add the new rows to the table
            console.log("populating new table data");
            var row;
            while (rows.length > 0) {
                row = rows.shift();
                tableBody.append(row);
            }
        }
    } else {
        if (panelDesc.name === 'trade')
            html = '<tr><td>nothing here...</td><td></td><td></td><td></td><td></td><td></td><td></td><td></td><td></td><td></td></tr>';
        else if (panelDesc.name === 'audit')
            html = '<tr><td>nothing here...</td><td></td><td></td><td></td><td></td><td></td><td></td><td></td><td></td></tr>'; // No action column
        $(panelDesc.tableID).html(html);
    }
}

// =================================================================================
//	Helpers for the filtering of trades
// =================================================================================
var filter = {};

/**
 * Describes all the fields that describe a trade.  Used to create
 * a filter that can be used to control which trades get shown in the
 * table.
 * @type {string[]}
 */
var names = [
    "cusip",
    "ticker",
    "par",
    "qty",
    "discount",
    "maturity",
    "issuer",
    "owner",
    "company"
];

/**
 * Parses the filter forms in the UI into an object for filtering
 * which trades are displayed in the table.
 * @param panelDesc An object describing which panel
 */
function processFilterForm(panelDesc) {
    "use strict";

    var form = document.forms[panelDesc.formID];

    console.log("Processing filter form");

    console.log(form.getElementsByTagName("input"));

    // Reset the filter parameters
    filter = {};

    // Build the filter based on the form inputs
    for (var i in names) {

        // Input ID example: "trade_owner"
        var name = names[i];
        var id = panelDesc.filterPrefix + name;

        if (form[id] && form[id].value !== "") {
            filter[name] = form[id].value;
        }
    }

    console.log("New filter parameters: " + JSON.stringify(filter));
    console.log("Rebuilding paper list");
    build_trades(bag.papers, panelDesc);
}

/**
 * Validates a trade object against a given set of filters.
 * @param entry The object to be validated.
 * @param filter The filter object to validate the trade against.
 * @returns {boolean} True if the trade is valid according to the filter, false otherwise.
 */
function excluded(entry, filter) {
    "use strict";

    if (filter.owner && filter.owner !== "" && entry.owner.toUpperCase().indexOf(filter.owner.toUpperCase()) == -1) return false;

    if (filter.issuer && filter.issuer !== "" && entry.issuer.toUpperCase().indexOf(filter.issuer.toUpperCase()) == -1) return false;

    if (filter.ticker && filter.ticker !== "" && entry.ticker.toUpperCase().indexOf(filter.ticker.toUpperCase()) == -1) return false;

    // Must be a valid trade if we reach this point
    return true;
}
