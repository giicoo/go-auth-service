function toggleForms(formId) {
    document.getElementById("registration-form").classList.add("hidden");
    document.getElementById("login-form").classList.add("hidden");
    document.getElementById("sessions").classList.add("hidden");
    document.getElementById("registration-error").classList.add("hidden");

    document.getElementById(formId).classList.remove("hidden");
  }

  function showSessions() {
    toggleForms("sessions");
  }

  // Handle registration form submission via Ajax
  document
    .getElementById("register-form")
    .addEventListener("submit", function (event) {
      event.preventDefault(); // Prevent page reload

      const email = document.getElementById("email").value;
      const password = document.getElementById("password").value;

      const data = {
        email: email,
        password: password,
      };

      // Send data to server via Fetch API (Ajax request)
      fetch("http://localhost:8080/create-user", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      })
        .then((response) => {
          if (response.ok) {
            return response.json(); // If successful, return JSON data
          } else {
            throw new Error(
              "Request failed with status " + response.status
            ); // If error, throw an exception
          }
        })
        .then((data) => {
          console.log("Registration successful", data);
          alert("Registration successful! Please log in.");
          toggleForms("login-form"); // Switch to login form
        })
        .catch((error) => {
          console.error("Error during registration:", error);
          document
            .getElementById("registration-error")
            .classList.remove("hidden"); // Show error message
        });
    });

    // Handle registration form submission via Ajax
document
.getElementById("login-form")
.addEventListener("submit", function (event) {
  event.preventDefault(); // Prevent page reload

  const email = document.getElementById("login-email").value;
  const password = document.getElementById("login-password").value;

  // Get the User-Agent string
  const userAgent = navigator.userAgent;

  // Get the User IP via a third-party API (ipinfo.io)
  fetch('https://ipinfo.io/json?token=f0e2261363429a') // Replace with your token from ipinfo.io
    .then(response => response.json())
    .then(ipData => {
      const userIp = ipData.ip; // User's IP address

      const data = {
        email: email,
        password: password,
        user_agent: userAgent, // Adding User-Agent to the data
        user_ip: userIp,       // Adding User-IP to the data
      };

      // Send data to server via Fetch API (Ajax request)
      fetch("http://localhost:8080/create-session", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      })
        .then((response) => {
          if (response.ok) {
            return response.json(); // If successful, return JSON data
          } else {
            throw new Error("Request failed with status " + response.status); // If error, throw an exception
          }
        })
        .then((data) => {
          console.log("Registration successful", data);
          alert("Registration successful! Please log in.");
          toggleForms("login-form"); // Switch to login form
        })
        .catch((error) => {
          console.error("Error during registration:", error);
          document
            .getElementById("login-error")
            .classList.remove("hidden"); // Show error message
        });
    })
    .catch((error) => {
      console.error("Failed to get IP address:", error);
      document
        .getElementById("login-error")
        .classList.remove("hidden"); // Show error message if IP fetch fails
    });
});