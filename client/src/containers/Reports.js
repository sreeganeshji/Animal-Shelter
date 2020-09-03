import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";

import { isLoggedIn } from "../services/auth";
import { NavLinks } from "../components";

const Reports = () => {
    const [state, setState] = useState({
        cookie: {}
    });

    // check auth
    useEffect(() => {
        const cookie = isLoggedIn();
        if (cookie.hasOwnProperty("isTrustedVolunteer") && cookie.isTrustedVolunteer) {
            setState({ ...state, cookie });
        }
        /* eslint-disable react-hooks/exhaustive-deps */
    }, []);

    return (
        <>
            <NavLinks isTrustedVolunteer={state.cookie.isTrustedVolunteer} />
            <main className="container">
                <h1>Reports</h1>
                <div className="interior-block">
                    <div className="form-unit">
                        <Link to="/monthly-adoption-report">Monthly Adoption Report</Link>
                    </div>
                    <div className="form-unit">
                        <Link to="/animal-control-report">Animal Control Report</Link>
                    </div>
                    <div className="form-unit">
                        <Link to="/expense-analysis">Expense Analysis</Link>
                    </div>
                    <div className="form-unit">
                        <Link to="/volunteer-lookup">Volunteer Lookup</Link>
                    </div>
                </div>
            </main>
        </>
    );
}

export default Reports;