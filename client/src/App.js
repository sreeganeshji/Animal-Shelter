import React from "react";
import {
  BrowserRouter as Router,
  Switch,
  Route
} from "react-router-dom";

import {
  AddAdoption,
  AddApplication,
  AddDog,
  AnimalControlReport,
  AnimalControlDrillDownReport,
  DogDashboard,
  DogDetails,
  ExpenseAnalysis,
  Login,
  MonthlyAdoptionReport,
  Reports,
  ReviewApplications,
  VolunteerLookup
} from "./containers";
import "./App.css";

export default function App() {
  return (
    <Router>
      <Routes />
    </Router>
  );
}

const Routes = () => {
  return (
    <Switch>
      <Route path="/add-adoption">
        <AddAdoption />
      </Route>
      <Route path="/add-application">
        <AddApplication />
      </Route>
      <Route path="/add-dog">
        <AddDog />
      </Route>
      <Route path="/animal-control-report/:monthAndYear">
        <AnimalControlDrillDownReport />
      </Route>
      <Route path="/animal-control-report">
        <AnimalControlReport />
      </Route>
      <Route path="/dog-dashboard">
        <DogDashboard />
      </Route>
      <Route path="/dog-details/:dogId">
        <DogDetails />
      </Route>
      <Route path="/expense-analysis">
        <ExpenseAnalysis />
      </Route>
      <Route path="/monthly-adoption-report">
        <MonthlyAdoptionReport />
      </Route>
      <Route path="/reports">
        <Reports />
      </Route>
      <Route path="/review-applications">
        <ReviewApplications />
      </Route>
      <Route path="/volunteer-lookup">
        <VolunteerLookup />
      </Route>
      <Route path="/">
        <Login />
      </Route>
    </Switch>
  )
}