import React from 'react';

const Footer = () => {
  return (
    // The main container with a dark background and consistent padding
    <footer className='w-full bg-black text-gray-400 py-12 px-4 sm:px-12 md:px-20'>
      <div className='max-w-6xl mx-auto'>
        {/* Contact Info (The "Questions? Contact us." line) */}
        <p className='text-base mb-6'>
          Questions? Call <a href='tel:0800-420-6969' className='hover:underline'>0800-420-6969</a> (Totally Fake Number)
        </p>

        {/* Grid Container for Links */}
        <div className='grid grid-cols-2 md:grid-cols-4 gap-4 text-sm'>
          {/* Column 1 */}
          <ul>
            <li className='mb-2 hover:underline cursor-pointer'>
              <a href='#'>FAQ</a>
            </li>
            <li className='mb-2 hover:underline cursor-pointer'>
              <a href='#'>Investor Relations</a>
            </li>
            <li className='mb-2 hover:underline cursor-pointer'>
              <a href='#'>Ways to Watch</a>
            </li>
            <li className='mb-2 hover:underline cursor-pointer'>
              <a href='#'>Fake Corporate Info</a>
            </li>
          </ul>

          {/* Column 2 */}
          <ul>
            <li className='mb-2 hover:underline cursor-pointer'>
              <a href='#'>Help Center</a>
            </li>
            <li className='mb-2 hover:underline cursor-pointer'>
              <a href='#'>Jobs (We're Hiring!)</a>
            </li>
            <li className='mb-2 hover:underline cursor-pointer'>
              <a href='#'>Terms of Use</a>
            </li>
            <li className='mb-2 hover:underline cursor-pointer'>
              <a href='#'>Contact Us</a>
            </li>
          </ul>

          {/* Column 3 */}
          <ul>
            <li className='mb-2 hover:underline cursor-pointer'>
              <a href='#'>Account</a>
            </li>
            <li className='mb-2 hover:underline cursor-pointer'>
              <a href='#'>Redeem Gift Cards</a>
            </li>
            <li className='mb-2 hover:underline cursor-pointer'>
              <a href='#'>Privacy Statement</a>
            </li>
            <li className='mb-2 hover:underline cursor-pointer'>
              <a href='#'>Fake Speed Test</a>
            </li>
          </ul>

          {/* Column 4 */}
          <ul>
            <li className='mb-2 hover:underline cursor-pointer'>
              <a href='#'>Media Center</a>
            </li>
            <li className='mb-2 hover:underline cursor-pointer'>
              <a href='#'>Buy Gift Cards</a>
            </li>
            <li className='mb-2 hover:underline cursor-pointer'>
              <a href='#'>Cookie Preferences</a>
            </li>
            <li className='mb-2 hover:underline cursor-pointer'>
              <a href='#'>Fake Legal Notices</a>
            </li>
          </ul>
        </div>

        {/* Service Code and Copyright/Region */}
        <div className='mt-8'>
          <button className='border border-gray-400 text-gray-400 text-xs px-3 py-1 mb-6 hover:text-white hover:border-white transition duration-200'>
            Service Code: FAKE-123-ABC
          </button>
          
          <p className='text-xs'>
            &copy; 2024 Fakeflix, Inc. (A Netflix Clone Project)
          </p>
          <p className='text-xs mt-1'>
            <span className='hover:underline cursor-pointer'>Faketown, State of Development.</span>
          </p>
        </div>
      </div>
    </footer>
  );
};

export default Footer;