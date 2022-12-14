import matplotlib.pyplot as plt

def make_plot(filename):
    data = []
    with open(filename, 'r') as f:
        data = f.read().splitlines()
    
    for x in data:
        print(x)
    
    data = [float(x) for x in data]

    plt.plot([a for a in range(len(data))], data)

if __name__ == '__main__':
    make_plot('../output_up')
    make_plot('../output_down')

    plt.show()
