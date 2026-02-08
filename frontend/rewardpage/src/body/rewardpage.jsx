import { Disclosure, DisclosureButton, DisclosurePanel, Menu, MenuButton, MenuItem, MenuItems} from '@headlessui/react';
import { Bars3Icon, BellIcon, XMarkIcon } from '@heroicons/react/24/outline';

import { useAuth } from '../hooks/useAuth';
import { useNavigate } from 'react-router-dom';

// import Example from './example';
import EarnPoints from './earnPoint';
import TaskDashboard from './TaskDashboard';
import ReferPoint from './referPoints';


const navigation = [
  { name: 'Dashboard', href: '#', current: true },
  { name: 'Team', href: '#', current: false },
  { name: 'Top Earner', href: '#', current: false },
  { name: 'Calendar', href: '#', current: false },
  { name: 'Notification', href: '#', current: false },
  { name: 'Sign out', href: '#', current: false },  
]

const userNavigation = [
  { name: 'Task for today', href: '#' },
  { name: 'Claim 100 gold', href: '#' },
  { name: 'Remember to Set your passwords', href: '#' },
]

function classNames(...classes) {
  return classes.filter(Boolean).join(' ')
}
export default function RewardPage(){
    const { logout } = useAuth();
    const navigate = useNavigate();

    const handleLogout = () => {
        logout();
        navigate('/login');
    };

    return (
        <>
        <div className="min-h-full">
        <Disclosure as= "nav" className="bg-white">       
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="flex h-16 items-center justify-between">
            
          <h2 className=" px-4 py-6 sm:px-4 lg:px-8 text-2xl font-medium text-gray-900">
           Rewards Hub </h2>        

        
        <div className="hidden md:block">
        <div className="ml-4 flex items-center md:ml-10">
            {/* Notification btn */}
            <Menu as="div" className="relative ml-3" >
            <MenuButton className="relative flex max-w-xs items-center rounded-full focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-500">
                <span className="absolute -inset-1.5" />
                <span className="sr-only">Open Notification menu</span>
            <button
             type="button"
             className="relative rounded-full border border-gray-400 p-1 text-gray-400 hover:text-black focus:outline-2 focus:outline-offset-2 focus:outline-black"
            >
                <span className="absolute -inset-1.5" />
                <span className="sr-only">View notifications</span>
                <BellIcon aria-hidden="true" className="size-6 " />
            </button>
            </MenuButton>
            
             <MenuItems
             transition
             className="absolute right-0 z-10 mt-2 w-48 origin-top-right rounded-md bg-white py-1 shadow-lg outline-1 outline-black/5 transition data-closed:scale-95 data-closed:transform data-closed:opacity-0 data-enter:duration-100 data-enter:ease-out data-leave:duration-75 data-leave:ease-in"
             >
                {userNavigation.map((item) => (
                    <MenuItem key={item.name}>
                    <a 
                    href={item.href}
                    className="block px-4 py-2 text-sm text-gray-700 data-focus:bg-gray-100 data-focus:outline-hidden"
                    > {item.name}
                    </a>
                    </MenuItem>
                ))}
             </MenuItems>
             </Menu>
            {/* Login btn */}
            <div className="relative ml-3">
          <button
            onClick={handleLogout}
            className=" bg-gray-500 text-white py-2 px-4 rounded-full hover:bg-black"
          >
            Sign out
          </button>
          </div>
          </div>
          </div>
          <div className="-mr-2 md:hidden">
            <DisclosureButton className="group relative inline-flex items-center justify-center rounded-md p-2 text-gray-400 hover:bg-white/5 hover:text-white focus:outline-2 focus:outline-offset-2 focus:outline-indigo-500">
            <span className="absolute -inset-0.5" />
                  <span className="sr-only">Open main menu</span>
                  <Bars3Icon aria-hidden="true" className="block size-6 group-data-open:hidden" />
                  <XMarkIcon aria-hidden="true" className="hidden size-6 group-data-open:block" />
            </DisclosureButton>
          </div>

        </div>  
        </div>

        <DisclosurePanel className="md: block">
            <div className="space-y-1 px-2 pt-2 pb-3 sm:px-3">
                {navigation.map((item)=>(
                    <DisclosureButton
                    key={item.name}
                    as= "a"
                    href={item.href}
                    aria-current = {item.current ? 'page': undefined}
                     className={classNames(
                    item.current ? 'bg-gray-900 text-white' : 'text-gray-300 hover:bg-white/5 hover:text-white',
                    'block rounded-md px-3 py-2 text-base font-medium',
                  )}                    
                    >
                    {item.name}
                    </DisclosureButton>
                ))}
            </div>
        </DisclosurePanel>
        </Disclosure>    

        <div className="mx-auto max-w-7xl sm:px-6 lg:px-8">
        <div className="flex h-5 items-center justify-between">
        <h2 className=" px-4 py-6 sm:px-4 lg:px-8 text-1 font-small text-gray-900">
           Earn point, unlock reward, and celebrate your progress </h2> 
        </div></div>

        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="flex h-16 items-center justify-between">
            <div className="ml-10 flex items-baseline space-x-2">
            
            <a
                href="#"
                className="box-border rounded-md bg-purple-300 px-3.5 py-2 text-sm font-semibold text-purple-800 shadow-xs
                        border-b-2 border-transparent
                        hover:border-purple-800 hover:bg-purple-200
                        transition-all duration-200"
            >
                Earn Points
            </a>

            <a
                href="#"
                className="box-border rounded-md px-3.5 py-2 text-sm font-semibold text-purple-800 shadow-xs
                        border-b-2 border-transparent
                        hover:border-purple-800 hover:bg-purple-200
                        transition-all duration-200"
            >
                Redeem Rewards
            </a>

            
        </div>
        </div>
          
        
        {/* TaskDashboard */}
        <TaskDashboard />

        {/* Earn Points page */}
        <EarnPoints/>

        {/* Refer Points page */}
        <ReferPoint />
        
        </div>

        

       
        </div>          
        </>
    )   
}