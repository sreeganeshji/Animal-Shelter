import React, { useEffect, useState } from "react";
import { Redirect } from "react-router-dom";

import { getData } from "../services/api";
import { isLoggedIn } from "../services/auth";
import { NavLinks } from "../components";

const AnimalControlReport = () => {

    const [state, setState] = useState({
        data: [],
        errorMessage: "",
        cookie: {},
        drillDownRedirect: false
    });

    // check auth
    useEffect(() => {
        const cookie = isLoggedIn();
        if (cookie.hasOwnProperty("isTrustedVolunteer") && cookie.isTrustedVolunteer) {
            setState({ ...state, cookie });
        }
        /* eslint-disable react-hooks/exhaustive-deps */
    }, []);

    useEffect(() => {
        if (state.cookie.hasOwnProperty("isTrustedVolunteer") && state.cookie.isTrustedVolunteer) {
            getData("/animal-control-report")
                .then(({ target }) => {
                    if (target.status === 200) {
                        setState({ ...state, data: JSON.parse(target.response) });
                    } else {
                        setState({ ...state, errorMessage: target.response });
                    }
                })
                .catch(error => {
                    console.log(error);
                    setState({ ...state, errorMessage: JSON.stringify(error) });
                })
        }
    }, [state.cookie]);
    return (
        <>
            {state.drillDownRedirect && <Redirect to={{ pathname: state.drillDownRedirect, state: { from: "/animal-control-report" } }} />}
            <NavLinks isTrustedVolunteer={state.cookie.isTrustedVolunteer} />
            <main className="container">
                <h1>Animal Control Report</h1>
                <table>
                    <thead>
                        <td>Month</td>
                        <td>Year</td>
                        <td>Total Dogs Count</td>
                        <td>Dogs 60 Days Count</td>
                        <td>Expenses</td>
                    </thead>
                    <tbody>
                        {state.data.length && (
                            state.data.map(datum => {
                                const expenseDollars = Math.floor(datum.expenses);
                                const expenseCents = (() => {
                                    const cents = Math.round((datum.expenses % 1) * 100);
                                    return cents < 10 ? "0" + cents : cents;
                                })();

                                return (
                                    <tr onClick={() => setState({ ...state, drillDownRedirect: `/animal-control-report/${datum.month.trim()}-${datum.year}` })}>
                                        <td>{datum.month}</td>
                                        <td>{datum.year}</td>
                                        <td>{datum.dogsTotalCount}</td>
                                        <td>{datum.dogsSixtyDaysCount}</td>
                                        <td>${expenseDollars}.{expenseCents}</td>
                                    </tr>
                                )
                            })
                        )}
                    </tbody>
                </table>

                {
                    state.errorMessage && (
                        <div styles={{ color: "red" }}>{state.errorMessage}</div>
                    )
                }
            </main>
        </>
    );
}

export default AnimalControlReport;