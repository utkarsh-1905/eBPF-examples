## What is what

So, in this program, we are basically linking out eBPF program to a network adapter on the Physical Layer with XDA.

Now the kernel space program is written in `ip.c` file.
It declares a buffer map, and pushes the recived data from the kernel into it. This shows that we can manipulate the traffic directly at the hardware level using the kernel without needing of programs in user space, and to add many such capabilities, we dont need to write a whole series of patches to the kernel.

The user space program is in `main.go` file which loads the eBPF functionality into the kernel, and reads data from the shared buffer and displays to the user.

Here the program is attached to the network driver, and gets triggered whenever a new packet is recieved. 

Upon recieving the packet, it gets the layer3 header (since ip is in layer3) after separating the physical layers using the offset. It the provides the header to `iphdr` struct to model the data and then pushes it to the buffer.

On user space side, every second, the program reads from the buffer, the latest value into a `[]byte` and spits out the info and console.