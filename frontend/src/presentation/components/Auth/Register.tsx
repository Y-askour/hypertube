import { useState } from 'react';
import { useTranslation } from 'react-i18next'
import { IoIosArrowForward } from 'react-icons/io'
import SignUp from './signup/Signup';

const Register = () => {
  const {t} = useTranslation();
  const [open,setOpen]=useState<boolean>(false);
  return (
    <>
        <p className='text-white'>
              {t('description')}
        </p>
        <div className='flex gap-2'>
            <input className={!open ? `px-3 py-4 rounded-full bg-white` : ''} placeholder={t('Email address')}/>
            <button onClick={()=>setOpen(true)} 
            className={!open ?  `px-3 py-4 bg-red-500 text-white font-bold rounded-full flex justify-center  items-center gap-5 ` : ''}>
                <span> {t('Get started')}</span>
                <IoIosArrowForward/>
            </button>
            {
              open && <SignUp setOpen={setOpen}/>
            }
        </div>
    </>
  )
}

export default Register