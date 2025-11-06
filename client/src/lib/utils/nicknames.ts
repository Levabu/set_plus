
const ADJECTIVES = ["Brave", "Witty", "Sneaky", "Loyal", "Bold", "Clever"];
const ANIMALS = ["Fox", "Penguin", "Tiger", "Dolphin", "Hawk", "Otter"];

export function generateNickname() {
  const adjective = ADJECTIVES[Math.floor(Math.random() * ADJECTIVES.length)];
  const animal = ANIMALS[Math.floor(Math.random() * ANIMALS.length)];
  return `${adjective}${animal}`;
}
