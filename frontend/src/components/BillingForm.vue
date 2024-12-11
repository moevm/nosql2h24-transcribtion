<template>
  <form @submit.prevent="handleBillingSubmit">
    <h2>Billing Information</h2>

    <div>
      <label>Card Number:</label>
      <input type="text" v-model="cardNumber" required placeholder="1234 5678 9012 3456" />
    </div>

    <div>
      <label>CVV:</label>
      <input type="text" v-model="cvv" required placeholder="123" maxlength="4" />
    </div>

    <div>
      <label>Expiration Date:</label>
      <input type="month" v-model="expiryDate" required />
    </div>

    <div>
      <label>Card Holder:</label>
      <input type="text" v-model="cardHolder" required placeholder="John Doe" />
    </div>

    <button type="submit">Submit</button>
    <button type="button" @click="$emit('close')">Cancel</button>
  </form>
</template>

<script>
import { ref } from 'vue';

export default {
  setup(_, { emit }) {
    const cardNumber = ref('');
    const cvv = ref('');
    const expiryDate = ref('');
    const cardHolder = ref('');

    const handleBillingSubmit = () => {
      console.log({
        cardNumber: cardNumber.value,
        cvv: cvv.value,
        expiryDate: expiryDate.value,
        cardHolder: cardHolder.value,
      });
      alert('Billing information submitted successfully!');
      emit('close');
    };

    return {
      cardNumber,
      cvv,
      expiryDate,
      cardHolder,
      handleBillingSubmit,
    };
  },
};
</script>

<style scoped>
form {
  display: flex;
  flex-direction: column;
  max-width: 400px;
  margin: auto;
  background-color: #f9f9f9;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

div {
  margin-bottom: 15px;
}

label {
  font-weight: bold;
  margin-bottom: 5px;
  display: block;
}

input {
  padding: 8px;
  width: 100%;
  box-sizing: border-box;
}

button {
  padding: 10px;
  margin-top: 10px;
  background-color: #42b983;
  color: white;
  border: none;
  cursor: pointer;
}

button[type="button"] {
  background-color: #f44336;
}

button:hover {
  background-color: #369870;
}
</style>
