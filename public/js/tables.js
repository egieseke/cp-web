/**
 * Created by davery on 3/29/2016.
 */
"use strict";
function createRow(data) {
    var tr = document.createElement('tr');

    for (var index in data) {
        var td = document.createElement('td');
        tr.appendChild(td);

        var text = document.createTextNode(data[index]);
        td.appendChild(text);
    }

    return tr;
}

/**
 * Generates a buy button cell that users can click to purchase commercial paper.
 * @param disabled True if the button should be disabled, false otherwise.
 * @param cusip The cusip for the paper that this button is assigned to.
 * @param issuer The issuer of the paper that this button is assigned to.
 * @returns {Element} A table cell with a configured buy button.
 */
function buyButton(disabled, vin, issuer) {
    var button = document.createElement('button');
    button.setAttribute('type', 'button');
    button.setAttribute('data_vin', vin);
    button.setAttribute('data_issuer', issuer);
    if(disabled) button.disabled = true;
    button.classList.add('buyPaper');
    button.classList.add('altButton');

    var span = document.createElement('span');
    span.classList.add('fa');
    span.classList.add('fa-exchange');
    span.innerHTML = ' &nbsp;&nbsp;Transfer Title';
    button.appendChild(span);

    // Wrap the buy button in a td like the other items in the row.
    var td = document.createElement('td');
    td.appendChild(button);

    return td;
}

function simulateViolationButton(disabled, licenseId, driver) {
    var button = document.createElement('button');
    button.setAttribute('type', 'button');
    button.setAttribute('data_licenseId', licenseId);
    button.setAttribute('data_driver', driver);
    if(disabled) button.disabled = true;
    button.classList.add('simulateViolation');
    button.classList.add('altButton');

    var span = document.createElement('span');
    span.classList.add('fa');
    span.classList.add('fa-exchange');
    span.innerHTML = ' &nbsp;&nbsp;Violation';
    button.appendChild(span);

    // Wrap the buy button in a td like the other items in the row.
    var td = document.createElement('td');
    td.appendChild(button);

    return td;
}

function simulateTollButton(disabled, registrationId, owner) {
    var button = document.createElement('button');
    button.setAttribute('type', 'button');
    button.setAttribute('data_registrationId', registrationId);
    button.setAttribute('data_owner', owner);
    if(disabled) button.disabled = true;
    button.classList.add('simulateToll');
    button.classList.add('altButton');

    var span = document.createElement('span');
    span.classList.add('fa');
    span.classList.add('fa-exchange');
    span.innerHTML = ' &nbsp;&nbsp;Toll';
    button.appendChild(span);

    // Wrap the buy button in a td like the other items in the row.
    var td = document.createElement('td');
    td.appendChild(button);

    return td;
}

function renewLicenseButton(disabled, licenseId, driver) {
    var button = document.createElement('button');
    button.setAttribute('type', 'button');
    button.setAttribute('data_licenseId', licenseId);
    button.setAttribute('data_driver', driver);
    if(disabled) button.disabled = true;
    button.classList.add('renewLicense');
    button.classList.add('altButton');

    var span = document.createElement('span');
    span.classList.add('fa');
    span.classList.add('fa-exchange');
    span.innerHTML = ' &nbsp;&nbsp;Renew License';
    button.appendChild(span);

    // Wrap the buy button in a td like the other items in the row.
    var td = document.createElement('td');
    td.appendChild(button);

    return td;
}

function simulateNotifyButton(disabled, entry) {
    var button = document.createElement('button');
    button.setAttribute('type', 'button');
    button.setAttribute('data_licenseId', entry.licenseId);
    button.setAttribute('data_owner', entry.owner);
    button.setAttribute('data_type', entry.type);
    button.setAttribute('data_id', entry.txId);
    button.setAttribute('data_amt', entry.tollAmt);
    button.setAttribute('data_loc', entry.location);
    button.setAttribute('data_date', entry.issueDate);

    if(disabled) button.disabled = true;
    button.classList.add('sendMail');
    button.classList.add('altButton');

    var span = document.createElement('span');
    span.classList.add('fa');
    span.classList.add('fa-exchange');
    span.innerHTML = ' &nbsp;&nbsp;Notify';
    button.appendChild(span);

    // Wrap the buy button in a td like the other items in the row.
    var td = document.createElement('td');
    td.appendChild(button);

    return td;
}

function renewRegistrationButton(disabled, registrationId, owner) {
    var button = document.createElement('button');
    button.setAttribute('type', 'button');
    button.setAttribute('data_registrationId', registrationId);
    button.setAttribute('data_owner', owner);
    if(disabled) button.disabled = true;
    button.classList.add('renewRegistration');
    button.classList.add('altButton');

    var span = document.createElement('span');
    span.classList.add('fa');
    span.classList.add('fa-exchange');
    span.innerHTML = ' &nbsp;&nbsp;Renew Registration';
    button.appendChild(span);

    // Wrap the buy button in a td like the other items in the row.
    var td = document.createElement('td');
    td.appendChild(button);

    return td;
}

function paper_to_entries(paper) {
    var entries = [];
        // Create a row for each valid trade
        var entry = {
            issueDate: paper.issueDate,
            vin: paper.vin,
            make: paper.make,
            model: paper.model,
            year: paper.year,
            color: paper.color,
            miles: paper.miles,
            value: paper.value,
            issuer: paper.issuer,
            owner: paper.owner
        };

        // Save which paper this is associated with
        entry.paper = paper;
        
        entries.push(entry);
    return entries;
}

function license_to_entries(license) {
    var entries = [];
        // Create a row for each valid license
        var entry = {
            licenseId: license.licenseId,
            testId: license.testId,
            address: license.address,
            city: license.city,
            state: license.state,
            zip: license.zip,
            driver: license.driver,
            auto: license.auto,
            issueDate: license.issueDate,
            expiryDate: license.expiryDate
        };

        // Save which paper this is associated with
        entry.license = license;

        entries.push(entry);
    return entries;
}

function registration_to_entries(registration) {
    var entries = [];
        var entry = {
            registrationId: registration.registrationId,
            plateNum: registration.plateNum,
            vin: registration.vin,
            testId: registration.testId,
            policyId: registration.policyId,
            owner: registration.owner,
            auto: registration.auto,
            issueDate: registration.issueDate,
            expiryDate: registration.expiryDate
        };

        // Save which paper this is associated with
        entry.registration= registration;

        entries.push(entry);
    return entries;
}

function toll_to_entries(toll) {
    var entries = [];
        var entry = {
            registrationId: toll.registrationId,
            txId: toll.txId,
            type: toll.type,
            //plateNum: toll.plateNum,
            owner: toll.owner,
            issueDate: toll.issueDate,
            tollAmt: toll.tollAmt,
	    location: toll.location
        };

        // Save which paper this is associated with
        entry.toll= toll;

        entries.push(entry);
    return entries;
}

function violation_to_entries(violation) {
    var entries = [];
        var entry = {
            registrationId: violation.registrationId,
            txId: violation.txId,
            licenseId: violation.licenseId,
            type: violation.type,
            //plateNum: violation.plateNum,
            owner: violation.driver,
            issueDate: violation.issueDate,
            tollAmt: violation.fine,
            location: violation.location
        };

        // Save which paper this is associated with
        entry.violation= violation;

        entries.push(entry);
    return entries;
}
