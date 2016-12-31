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
function buyButton(disabled, cusip, issuer) {
    var button = document.createElement('button');
    button.setAttribute('type', 'button');
    button.setAttribute('href', '#transferModal');
    button.setAttribute('data_cusip', cusip);
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

function renewLicenseButton(disabled, licenseId, driver) {
    var button = document.createElement('button');
    button.setAttribute('type', 'button');
    button.setAttribute('href', '#licenseModal');
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

function renewRegistrationButton(disabled, registrationId, owner) {
    var button = document.createElement('button');
    button.setAttribute('type', 'button');
    button.setAttribute('href', '#registrationModal');
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
    for (var owner in paper.owner) {
        // Create a row for each valid trade
        var entry = {
            issueDate: paper.issueDate,
            cusip: paper.cusip,
            ticker: paper.ticker,
            par: paper.par,
            quantity: paper.owner[owner].quantity,
            discount: paper.discount,
            maturity: paper.maturity,
            issuer: paper.issuer,
            owner: paper.owner[owner].company
        };

        // Save which paper this is associated with
        entry.paper = paper;
        
        entries.push(entry);
    }
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
            issueDate: registration.issueDate,
            expiryDate: registration.expiryDate
        };

        // Save which paper this is associated with
        entry.registration= registration;

        entries.push(entry);
    return entries;
}
