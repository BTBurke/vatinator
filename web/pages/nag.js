import React from "react";
import Nav from "../components/nav";

export default function NagPage() {
  const barstyle = {
    width: "14%",
  };

  return (
    <>
      <div className="container mx-auto">
        <Nav />
        <div className="w-full lg:w-3/4 mx-auto text-white bg-primary">
          <h1 className="text-2xl text-accent-1 lg:text-4xl font-bold py-2">
            RIP Vatinator 2020-2026
          </h1>
          <div className="py-1">
            Sorry friends, the time has come to shut the Vatinator down.
          </div>
        </div>
      </div>
    </>
  );
}
