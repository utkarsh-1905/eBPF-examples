//go:build ignore

#include <linux/bpf.h>
#include <linux/ip.h>
#include <linux/if_ether.h>
#include <bpf/bpf_helpers.h>

struct Info {
    __u32 source_ip;
    __u32 dest_ip;
    __u8 ttl;
    __u8 protocol;
};

struct {
    __uint(type, BPF_MAP_TYPE_ARRAY);
    __type(key, __u32);
    __type(value, struct Info);
    __uint(max_entries, 255);
} ips SEC(".maps");

SEC("xdp") 
int get_ips(struct xdp_md *ctx) {
    bpf_printk("got a packet\n");     
    void *data_end = (void *)(long)ctx->data_end;
    void *data     = (void *)(long)ctx->data;
    struct ethhdr *eth = data;

    // check packet size
    if (eth + 1 > data_end) {
        return XDP_PASS;
    }

    // get the source address of the packet
    struct iphdr *iph = data + sizeof(struct ethhdr);
    if (iph + 1 > data_end) {
        return XDP_PASS;
    }

    __u32 ip_src = iph->saddr;
    bpf_printk("source ip address is %s\n", ip_src);

    struct Info d;

    d.source_ip = ip_src;
    d.dest_ip = iph->daddr;
    d.protocol = iph->protocol;
    d.ttl = iph->ttl;

    __u32 key = 0;

    bpf_printk("starting xdp ip filter\n");
    bpf_map_update_elem(&ips, &key, &d, BPF_ANY);
    return XDP_PASS;
}

char __license[] SEC("license") = "Dual MIT/GPL";