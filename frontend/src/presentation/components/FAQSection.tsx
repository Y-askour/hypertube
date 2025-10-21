import React, { useState } from 'react';
import { Minus, Plus } from 'lucide-react'; // You'll need to install 'lucide-react' for the icons

// Individual FAQ item component (defined below)
const FAQItem = ({ question, answer }) => {
  const [isOpen, setIsOpen] = useState(false);

  const toggleOpen = () => {
    setIsOpen(!isOpen);
  };

  return (
    <div className="mb-2 bg-[#303030] max-w-4xl mx-auto">
      {/* Question Header (The clickable part) */}
      <button
        className="flex justify-between items-center w-full p-6 text-2xl font-normal text-white hover:bg-[#414141] transition-colors duration-200"
        onClick={toggleOpen}
        aria-expanded={isOpen}
      >
        <span>{question}</span>
        {/* Toggle Icon: Plus when closed, Minus when open */}
        {isOpen ? (
          <Minus className="w-8 h-8 text-white transform transition-transform duration-300" />
        ) : (
          <Plus className="w-8 h-8 text-white transform transition-transform duration-300" />
        )}
      </button>

      {/* Answer Content (The expandable child) */}
      {/* Uses a height/max-height transition for a smooth animation effect */}
      <div
        className={`overflow-hidden transition-all duration-300 ease-in-out ${
          isOpen ? 'max-h-96 opacity-100 py-6 border-t-2 border-black' : 'max-h-0 opacity-0'
        }`}
      >
        <p className="px-6 text-xl text-white whitespace-pre-wrap leading-relaxed">
          {answer}
        </p>
      </div>
    </div>
  );
};

// Main Section Component
const FAQSection = () => {
  const faqData = [
    {
      id: 1,
      question: 'What is Fakeflix?',
      answer:
        'Fakeflix is a fictional streaming service created as a web development project clone. It offers an imagined library of movies and TV shows across many genres. It is NOT the real Netflix and no actual content is streamed.',
    },
    {
      id: 2,
      question: 'How much does Fakeflix cost?',
      answer:
        'Since this is a CLONE project, it costs \$0.00! There are no subscription fees, and no financial data is ever required or collected. Enjoy the free development experience!',
    },
    {
      id: 3,
      question: 'Where can I watch?',
      answer:
        'You can watch the Fakeflix clone project right here in your web browser! It is fully responsive and works on phones, tablets, and desktops (as long as they support modern web standards).',
    },
    {
      id: 4,
      question: 'How do I cancel?',
      answer:
        "There's nothing to cancel! Just close the tab or navigate away. Since there is no account, subscription, or payment information, your 'cancellation' is instant and hassle-free.",
    },
    {
      id: 5,
      question: 'What can I watch on Fakeflix?',
      answer:
        'You can watch high-quality placeholder images, creative UI designs, and a demonstration of React and Tailwind CSS proficiency. Actual video content is not included.',
    },
    {
      id: 6,
      question: 'Is Fakeflix good for kids?',
      answer:
        'Yes, this project is designed to be safe for all ages as it contains no actual mature content, only mock data and UI elements.',
    },
  ];

  return (
    <div className='w-full bg-black py-12 md:py-20 border-b-8 border-[#222]'>
      <h2 className='text-white text-3xl md:text-5xl font-bold text-center mb-10'>
        Frequently Asked Questions
      </h2>
      <div className='max-w-4xl mx-auto px-4'>
        {faqData.map((item) => (
          <FAQItem key={item.id} question={item.question} answer={item.answer} />
        ))}
      </div>
      
      {/* Call to action section for the clone */}
      <div className="text-center mt-10 text-white text-xl">
          Ready to try the real thing? <br/>
          <button className="bg-red-600 hover:bg-red-700 text-white font-medium py-3 px-6 rounded mt-4 transition-colors duration-300">
              Try the Demo Now (Fake Button)
          </button>
      </div>
    </div>
  );
};

export default FAQSection;