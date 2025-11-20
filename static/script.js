// customElements.define(
//   "c-persist",
//   class extends HTMLElement {
//     connectedCallback() {
//       const inputs = this.querySelectorAll("input");
//       inputs.forEach(input => {
//         const key = input.name
//         let value = localStorage.getItem(key);
//         // if (!value) {
//         //   value = prompt("Digite seu nome");
//         // }
//         input.addEventListener("onchange", () => {
//           console.log("onchange");
//           if (value) {
//             localStorage.setItem(key, value);
//           }
//         });
//         if (value) {
//           input.value = value
//         }
//       });
//     }
//   }
// );
