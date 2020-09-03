import React, { useEffect, useState } from "react";

import { getData } from "../services/api";
import { isLoggedIn } from "../services/auth";
import { NavLinks } from "../components";

function getMonthName(monthNum) {
    switch (monthNum) {
        case 0: return "January"
        case 1: return "February"
        case 2: return "March"
        case 3: return "April"
        case 4: return "May"
        case 5: return "June"
        case 6: return "July"
        case 7: return "August"
        case 8: return "September"
        case 9: return "October"
        case 10: return "November"
        case 11: return "December"

        default: return "month"

    }
}

function GetMonthOffset(offset) {
    let today = new Date()
    today.setMonth(today.getMonth() + offset)

    return getMonthName(today.getMonth()) + " " + today.getFullYear()
}

const MonthlyAdoptionReport = () => {

    const [state, setState] = useState({
        data: {},
        acc:{},
        errorMessage: "",
        loginRedirect: false,
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

    useEffect(() => {
        for (let i = -11; i <= 0; i++) {
            UpdateTables(i)
        }

    }, [state.cookie]);

    function UpdateTables(ind) {

        if (state.cookie.hasOwnProperty("isTrustedVolunteer") && state.cookie.isTrustedVolunteer) {
            getData("/monthly-adoption-report?month=" + ind)
                .then(({ target }) => {
                    if (target.status === 200) {

                        const dataDict = state.data
                        dataDict[ind] = JSON.parse(target.response)

                        setState({ ...state, data: dataDict });

                    } else {
                        setState({ ...state, errorMessage: target.response });
                    }
                })
                .catch(error => {
                    console.log(error);
                    setState({ ...state, errorMessage: JSON.stringify(error) });
                })
        }
    }

    function getTotalValues(indx) {
        var acc = {"SurrenderCount":0.0, "AdoptionCount":0.0, "AdoptionFees":0.0, "Expenses":0.0, "Profit":0.0}

        state.data[indx].map(datum=> {
            acc.SurrenderCount += datum.SurrenderCount.Float64
            acc.AdoptionCount += datum.AdoptionCount.Float64
            acc.AdoptionFees += datum.AdoptionFees.Float64
            acc.Expenses += datum.Expenses.Float64
            acc.Profit += datum.Profit.Float64
        })

        const expenseCents = (() => {
            const cents = Math.round(acc.Expenses % 100);
            return cents < 10 ? "0" + cents : cents;
        })();
        const expenseDollars = Math.floor(acc.Expenses / 100);

        const adoptionFeeCents = (() => {
            const cents = Math.round(acc.AdoptionFees % 100);
            return cents < 10 ? "0" + cents : cents;
        })();
        const adoptionFeeDollars = Math.floor(acc.AdoptionFees / 100);

        const profitCents = (() => {
            const cents = Math.round(acc.Profit % 100);
            return cents < 10 ? "0" + cents : cents;
        })();
        const profitDollars = Math.floor(acc.Profit / 100);

        return (
            <>
                <td></td>
                <td>{acc.SurrenderCount}</td>
                <td>{acc.AdoptionCount}</td>
                <td>${expenseDollars}.{expenseCents}</td>
                <td>${adoptionFeeDollars}.{adoptionFeeCents}</td>
                <td>${profitDollars}.{profitCents}</td>
           </>
        )
    }

    function getTableValues() {
        var indices = [];
        var i = 0;
        for (i = -11; i <= 0; i++) {
            indices.push(i);
        }

        return (
            indices.map(indx => {
                if (indx in state.data && state.data[indx].length) {


                    return (
                        <>
                            <tr>

                                <td>
                                    <h3>{GetMonthOffset(indx)}</h3></td>
                                <td></td>
                                <td></td>
                                <td></td>
                                <td></td>
                                <td></td>
                                <td></td>
                            </tr>
                            {state.data[indx].map(datum => {
                                const expenseCents = (() => {
                                    const cents = Math.round(datum.Expenses.Float64 % 100);
                                    return cents < 10 ? "0" + cents : cents;
                                })();
                                const expenseDollars = Math.floor(datum.Expenses.Float64 / 100);

                                const adoptionFeeCents = (() => {
                                    const cents = Math.round(datum.AdoptionFees.Float64 % 100);
                                    return cents < 10 ? "0" + cents : cents;
                                })();
                                const adoptionFeeDollars = Math.floor(datum.AdoptionFees.Float64 / 100);

                                const profitCents = (() => {
                                    const cents = Math.round(datum.Profit.Float64 % 100);
                                    return cents < 10 ? "0" + cents : cents;
                                })();
                                const profitDollars = Math.floor(datum.Profit.Float64 / 100);

                                return (
                                    <tr>
                                    <td></td>
                                        <td>{datum.breed}</td>
                                        <td>{datum.SurrenderCount.Float64}</td>
                                        <td>{datum.AdoptionCount.Float64}</td>
                                        <td>${expenseDollars}.{expenseCents}</td>
                                        <td>${adoptionFeeDollars}.{adoptionFeeCents}</td>
                                        <td>${profitDollars}.{profitCents}</td>
                                    </tr>
                                )
                            }
                            )}
                            <tr>
                                <td><h3>Total</h3></td>
                                {
                                    getTotalValues(indx)
                                }
                            </tr>
                        </>

                    )

                }
            }

            )

        )
    }

    return (
        <>
            <NavLinks isTrustedVolunteer={state.cookie.isTrustedVolunteer} />
            <main className="container">
                <h1>Monthly Adoption Report</h1>
                <table>
                    <thead>
                        <td>Month</td>
                        <td>Breed</td>
                        <td>Surrender Count</td>
                        <td>Adoption Count</td>
                        <td>Expenses</td>
                        <td>Adoption Fees</td>
                        <td>Profit</td>
                    </thead>
                    <tbody>
                        {
                            getTableValues()
                        }
                    </tbody>
                </table>

                {state.errorMessage && (
                    <div styles={{ color: "red" }}>{state.errorMessage}</div>
                )}
            </main>
        </>
    );
}

export default MonthlyAdoptionReport;

